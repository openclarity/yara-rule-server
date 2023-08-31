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
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"

	"github.com/openclarity/yara-rule-server/pkg/config"
)

func DownloadAndCompile(cfg *config.Config, logger *logrus.Entry) error {
	archives := download(cfg.RulePath, cfg.RuleURLs, logger)
	if len(archives) == 0 {
		return errors.New("there is no successfully downloaded rules")
	}
	if num := unarchive(archives, logger); num == 0 {
		return errors.New("there is no successfully unarchived rules")
	}

	if err := generateIndex(cfg.IndexGenPath, logger); err != nil {
		return fmt.Errorf("failed to generate index.yar: %v", err)
	}
	if err := compile(cfg.YaracPath, "index.yar", cfg.RulePath, logger); err != nil {
		return fmt.Errorf("failed to copile index.yar: %v", err)
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
