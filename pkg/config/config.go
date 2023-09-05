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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	CacheDir = "/var/lib/yara-rule-server"
	RulePath = "/var/lib/yara-rule-server/compiled"
)

type Config struct {
	LogLevel                  string       `mapstructure:"log_level"`
	EnableJSONLog             bool         `mapstructure:"enable_json_log"`
	RuleSources               []RuleSource `mapstructure:"rule_sources"`
	YaracPath                 string       `mapstructure:"yarac_path"`
	IndexGenPath              string       `mapstructure:"index_gen_path"`
	RuleUpdateSchedule        string       `mapstructure:"rule_update_schedule"`
	ServerAddress             string       `mapstructure:"server_address"`
	HealthCheckAddressAddress string       `mapstructure:"health_check_address"`
}

type RuleSource struct {
	Name         string `mapstructure:"name"`
	URL          string `mapstructure:"url"`
	ExcludeRegex string `mapstructure:"exclude_regex"`
}

func LoadConfig(cfgFile string) *Config {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigType("yaml")
		viper.SetConfigFile("/etc/yara-rule-server/config.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match
	setDefaults()

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	// Load config
	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	cobra.CheckErr(err)

	return cfg
}

func setDefaults() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("yarac_path", "/usr/bin/yarac")
	viper.SetDefault("rule_update_schedule", "0 0 * * *")
	viper.SetDefault("server_address", ":8080")
	viper.SetDefault("health_check_address", ":8082")
}
