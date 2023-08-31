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
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func download(rulePath string, urls []string, logger *logrus.Entry) []string {
	archives := make([]string, 0)
	for i := range urls {
		u, err := url.Parse(urls[i])
		if err != nil {
			logger.Errorf("Failed to parse url=%s: %v", urls[i], err)
			continue
		}

		dirName := fmt.Sprintf("%s%s", path.Dir(rulePath), path.Dir(u.Path))
		if err := os.MkdirAll(dirName, 0755); err != nil {
			logger.Errorf("failed to create directory=%s: %v", dirName, err)
			continue
		}
		fileName := fmt.Sprintf("%s/%s", dirName, path.Base(u.Path))
		logger.Infof("Downloading %s into %s", urls[i], fileName)
		if err := downloadFile(fileName, urls[i]); err != nil {
			logger.Errorf("Failed to download file from url=%s: %v", urls[i], err)
			continue
		}
		archives = append(archives, fileName)
	}
	return archives
}

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create ouptut file=%s: %v", filepath, err)
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get url: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get url=%s: %s", url, resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file=%s: %v", filepath, err)
	}

	return nil
}
