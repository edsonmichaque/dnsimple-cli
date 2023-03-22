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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	profile    string
)

const (
	binaryName              = "dnsimple"
	configAccessToken       = "access-token"
	configAccount           = "account"
	configBaseURL           = "base-url"
	configCollaboratorID    = "collaborator-id"
	configConfigFile        = "config-file"
	configConfirm           = "confirm"
	configDomain            = "domain"
	defaultConfigFileFormat = "yaml"
	defaultProfile          = "default"
	envDNSimpleConfigFile   = "DNSIMPLE_CONFIG_FILE"
	envDNSimpleProfile      = "DNSIMPLE_PROFILE"
	envDev                  = "DEV"
	envPrefix               = "DNSIMPLE"
	envProd                 = "PROD"
	envSandbox              = "SANDBOX"
	envXDGConfigHome        = "XDG_CONFIG_HOME"
	flagOutput              = "output"
	flagProfile             = "profile"
	flagQuery               = "query"
	flagRecordID            = "record-id"
	flagSandbox             = "sandbox"
	formatJSON              = "json"
	formatTable             = "table"
	formatText              = "text"
	formatYAML              = "yaml"
	optPage                 = "page"
	optPerPage              = "per-page"
	optionFromFile          = "from-file"
	pathConfigFile          = "/etc/dnsimple"
	pathDNSimple            = "dnsimple"
)

func Run(opts *Options) error {
	return CmdRoot(opts).Execute()
}

func CmdRoot(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:          binaryName,
		SilenceUsage: true,
	}, opts)

	cmd.AddCommand(CmdAccounts(opts))
	cmd.AddCommand(CmdConfig(opts))
	cmd.AddCommand(CmdDomain(opts))
	cmd.AddCommand(CmdVersion(opts))
	cmd.AddCommand(CmdWhoami(opts))

	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().Bool(flagSandbox, false, "Sandbox environment")
	cmd.PersistentFlags().String(configAccessToken, "", "Access token")
	cmd.PersistentFlags().String(configAccount, "", "Account")
	cmd.PersistentFlags().String(configBaseURL, "", "Base URL")
	cmd.PersistentFlags().StringVar(&profile, flagProfile, defaultProfile, "Profile")
	cmd.PersistentFlags().StringVarP(&configFile, configConfigFile, "c", "", "Configuration file")

	cmd.MarkFlagsMutuallyExclusive(configBaseURL, flagSandbox)

	viper.SetEnvPrefix(envPrefix)
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		panic(err)
	}

	return cmd
}

func initConfig() {
	var err error

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else if configFile := os.Getenv(envDNSimpleConfigFile); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		configHome := os.Getenv(envXDGConfigHome)
		if configHome == "" {
			configHome, err = os.UserConfigDir()
			cobra.CheckErr(err)
		}

		viper.AddConfigPath(filepath.Join(configHome, pathDNSimple))
		viper.AddConfigPath(pathConfigFile)
		viper.SetConfigType(defaultConfigFileFormat)

		viper.SetConfigName(defaultProfile)

		if profileEnv := os.Getenv(envDNSimpleProfile); profileEnv != "" {
			viper.SetConfigName(profileEnv)
		}

		if profile != "" {
			viper.SetConfigName(profile)
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Found error: ", err.Error())
		}
	}
}

func createCmd(cmd *cobra.Command, opts *Options) *cobra.Command {
	applyOpts(cmd, opts)

	return cmd
}
