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

package config

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel      string `mapstructure:"LOG_LEVEL"`
	EnableJSONLog bool   `mapstructure:"ENABLE_JSON_LOG"`
}

func LoadConfig(cfgFile string) *Config {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kubeclarity" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".yara-rule-server")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	// Load config
	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	cobra.CheckErr(err)

	return cfg
}
