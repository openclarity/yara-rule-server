// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rules

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/yara-rule-server/pkg/config"
)

func DownloadAndCompile(cfg *config.Config, logger *logrus.Entry) error {
	// First try to download new copies of all the sources
	yarFilesToIndex := make([]string, 0)
	tempDir := path.Join(config.CacheDir, "tmp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create yara rule server temp directory. tempDir=%s: %v", tempDir, err)
	}
	for _, source := range cfg.RuleSources {
		// Create directory for this source if it doesn't exist
		sourceDir := path.Join(config.CacheDir, "sources", source.Name)
		if err := os.MkdirAll(sourceDir, 0755); err != nil {
			logger.Errorf("Failed to create directory %s. Using the last successful download if available: %v", sourceDir, err)
			continue
		}

		err := atomicDownloadAndReplace(source, sourceDir, tempDir, logger)
		if err != nil {
			logger.Errorf("unable to update source %s; using any existing data: %v", source.Name, err)
		}

		var reg *regexp.Regexp
		if source.ExcludeRegex != "" {
			var err error
			reg, err = regexp.Compile(source.ExcludeRegex)
			if err != nil {
				logger.Errorf("Failed to compile regexp %s: %v", source.ExcludeRegex, err)
				continue
			}
		}

		yarFiles, err := createYarFileListToIndex(sourceDir, reg)
		if err != nil {
			logger.Errorf("failed to create list of files from source %s for indexing: %v", source.Name, err)
			continue
		}
		yarFilesToIndex = append(yarFilesToIndex, yarFiles...)
	}

	if err := generateIndexAndCompile(cfg.YaracPath, yarFilesToIndex, tempDir); err != nil {
		return fmt.Errorf("failed to create compiled rules: %v", err)
	}

	return nil
}

func ScheduledDownload(cfg *config.Config, logger *logrus.Entry) (*gocron.Scheduler, error) {
	s := gocron.NewScheduler(time.UTC)

	// cron expressions supported
	if _, err := s.Cron(cfg.RuleUpdateSchedule).Do(DownloadAndCompile, cfg, logger); err != nil {
		return nil, fmt.Errorf("failed to create cron job: %v", err)
	}

	// starts the scheduler asynchronously
	s.StartAsync()

	return s, nil
}

// atomicDownloadAndReplace will download and unarchive the URL,
// only if those are successful will the data replace sourceDir
func atomicDownloadAndReplace(source config.RuleSource, sourceDir, tempDir string, logger *logrus.Entry) error {
	tmpSourceDir, err := os.MkdirTemp(tempDir, source.Name+"-yara-rule")
	if err != nil {
		return fmt.Errorf("failed to create temp directory for %s Using the last successful download if available: %v", source.Name, err)
	}

	fileName := filepath.Join(tmpSourceDir, source.Name+".zip")
	logger.Infof("Downloading %s into %s", source.URL, fileName)
	if err := downloadFile(fileName, source.URL); err != nil {
		return fmt.Errorf("failed to download rule source %s. Using the last successful download if available. URL=%s: %v", source.Name, source.URL, err)
	}

	logger.Infof("Unarchive rules file=%s into %s", fileName, tmpSourceDir)
	if err := unzip(fileName, tmpSourceDir); err != nil {
		return fmt.Errorf("failed to unarchive source %s. Using last successful download if available. File=%s: %v", source.Name, fileName, err)
	}

	// Replace contents of source dir with the downloaded and unarchived data
	if err := os.RemoveAll(sourceDir); err != nil {
		return fmt.Errorf("failed to remove previous sources: %v", err)
	}
	if err := os.Rename(tmpSourceDir, sourceDir); err != nil {
		return fmt.Errorf("failed to move downloaded source: %v", err)
	}

	return nil
}
