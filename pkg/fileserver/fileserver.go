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
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/openclarity/yara-rule-server/pkg/config"
)

func Start(cfg *config.Config, logger *logrus.Entry) *http.Server {
	logger.Infof("Starting file server. Rule file: %s", cfg.RulePath)
	sFile := func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, cfg.RulePath)
	}

	http.HandleFunc("/", sFile)
	server := &http.Server{Addr: cfg.ServerAddress, Handler: nil}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Failed to start http server: %v", err)
		}
	}()

	return server
}
