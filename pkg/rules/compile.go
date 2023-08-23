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
	"os/exec"

	"github.com/sirupsen/logrus"
)

func generateIndex(indexGenPATH string, logger *logrus.Entry) error {
	indexGen := exec.Command(indexGenPATH)

	resultJsonB, err := indexGen.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run index_gen.sh: %v, %v", err, string(resultJsonB))
	}
	logger.Debugln(string(resultJsonB))

	return nil
}

func compile(yaracPATH, input, output string, logger *logrus.Entry) error {
	yarac := exec.Command(yaracPATH, input, output)

	resultJsonB, err := yarac.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run yarac: %v, %v", err, string(resultJsonB))
	}
	logger.Debugln(string(resultJsonB))

	return nil
}
