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

package fileserver

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func Start(path string, logger *logrus.Entry) {
	logger.Infof("Rule file: %s", path)
	sFile := func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, path)
	}

	http.HandleFunc("/", sFile)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Errorf("Failed to start http server: %v", err)
	}
}
