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

package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openclarity/yara-rule-server/pkg/config"
	"github.com/openclarity/yara-rule-server/pkg/fileserver"
	"github.com/openclarity/yara-rule-server/pkg/rules"
)

var (
	cfgFile string
	cfg     *config.Config
	logger  *logrus.Entry
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yara-rule-server",
	Short: "Yara rule server",
	Long:  `Yara rule server downloads Yara rules, compiles them, and provides compiled rules as a file server`,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		initConfig,
		initLogger,
	)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yara-rules-server.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg = config.LoadConfig(cfgFile)
	fmt.Printf("Configuration: %+v", cfg)
}

func initLogger() {
	log := logrus.New()
	if level, err := logrus.ParseLevel(cfg.LogLevel); err != nil {
		log.SetLevel(level)
	}
	if cfg.EnableJSONLog {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	logger = log.WithField("app", "yara-rule-server")
}

func run(cmd *cobra.Command, args []string) {
	archives := rules.Download(cfg.RuleURLs, logger)
	rules.Unarchive(archives, logger)
	if err := rules.Compile(cfg.RulePath); err != nil {
		logger.Fatalf("Falied to compile YARA rules: %v", err)
	}
	fileserver.Start(cfg.RulePath, logger)
}
