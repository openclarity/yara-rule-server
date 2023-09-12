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
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/openclarity/yara-rule-server/pkg/config"
)

func createYarFileListToIndex(sourceDir string, reg *regexp.Regexp) ([]string, error) {
	yarFilesToIndex := make([]string, 0)
	err := filepath.WalkDir(sourceDir, func(path string, file fs.DirEntry, err error) error {
		// If there was an error walking this file structure return it
		if err != nil {
			return err
		}

		// We are only looking for files ending in .yar.
		// So return if it does not have that suffix or is a directory.
		if !strings.HasSuffix(path, ".yar") || file.IsDir() {
			return nil
		}

		// We've got a yar file, now we need to check that its not one
		// that should be excluded based on the configuration.
		// If we match the regex return nil so that this file is ignored.
		if reg != nil && reg.MatchString(path) {
			return nil
		}

		yarFilesToIndex = append(yarFilesToIndex, path)

		return nil
	})

	return yarFilesToIndex, err
}

func generateIndexAndCompile(yaracPATH string, yarFilesToIndex []string, tempDir string) error {
	// Write index file so that we can pass it to yarac
	tmpIndexFile, err := os.CreateTemp(tempDir, "index")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}

	for _, yarFile := range yarFilesToIndex {
		tmpIndexFile.WriteString(fmt.Sprintf("include \"%s\"\n", yarFile)) // nolint:errcheck
	}
	tmpIndexFile.Close()

	// Generate compiled rules in a temp directory
	tmpCompiledFile, err := os.CreateTemp(tempDir, "compiled")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	err = compile(yaracPATH, tmpIndexFile.Name(), tmpCompiledFile.Name())
	if err != nil {
		return fmt.Errorf("failed to compile %s: %v", tmpIndexFile.Name(), err)
	}

	// Remove index temp index file
	os.Remove(tmpIndexFile.Name())

	// Now we have the compile rules atomically move
	// it into the location to be served by the http server.
	if err := os.Rename(tmpCompiledFile.Name(), config.RulePath); err != nil {
		return fmt.Errorf("failed to move compiled yara rules: %v", err)
	}

	return nil
}

func compile(yaracPATH, input, output string) error {
	yarac := exec.Command(yaracPATH, "-w", input, output)

	stdoutStderrBytes, err := yarac.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run yarac: %v, %v", err, string(stdoutStderrBytes))
	}

	return nil
}
