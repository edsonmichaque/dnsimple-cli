// Copyright 2023 Edson Michaque
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
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/build"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	profile    string
)

const (
	flagAccount             = "account"
	flagAccessToken         = "access-token"
	flagBaseURL             = "base-url"
	flagSandbox             = "sandbox"
	flagProfile             = "profile"
	flagConfig              = "config-file"
	envPrefix               = "DNSIMPLE"
	defaultProfile          = "default"
	defaultConfigFileFormat = "yaml"
)

func NewCmdRoot(opts *internal.CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "dnsimple",
		Version:      build.Version,
		SilenceUsage: true,
	}

	cobra.OnInitialize(lookupConfigFiles)

	cmd.PersistentFlags().String(flagAccount, "", "Account")
	cmd.PersistentFlags().String(flagBaseURL, "", "Base URL")
	cmd.PersistentFlags().String(flagAccessToken, "", "Access token")
	cmd.PersistentFlags().Bool(flagSandbox, false, "Sandbox environment")
	cmd.PersistentFlags().StringVarP(&configFile, flagConfig, "c", "", "Configuration file")
	cmd.PersistentFlags().StringVar(&profile, flagProfile, "", "Profile")

	cmd.MarkFlagsMutuallyExclusive(flagBaseURL, flagSandbox)

	viper.SetEnvPrefix(envPrefix)
	_ = viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}

func lookupConfigFiles() {
	var err error

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else if configFile := os.Getenv("DNSIMPLE_CONFIG_FILE"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			configHome, err = os.UserConfigDir()
			cobra.CheckErr(err)
		}

		viper.AddConfigPath(filepath.Join(configHome, "dnsimple"))
		viper.AddConfigPath("/etc/dnsimple")
		viper.SetConfigType(defaultConfigFileFormat)

		if profile != "" {
			viper.SetConfigName(profile)
		} else if configProfile := os.Getenv("DNSIMPLE_PROFILE"); configProfile != "" {
			viper.SetConfigName(configProfile)
		} else {
			viper.SetConfigName(defaultProfile)
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Found error: ", err.Error())
		}
	}
}
