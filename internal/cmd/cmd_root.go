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
	flagConfig              = "config"
	envPrefix               = "DNSIMPLE"
	defaultProfile          = "default"
	defaultConfigFileFormat = "yaml"
)

func NewCmdRoot(opts *internal.CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dnsimple",
		Version: build.Version,
	}

	cobra.OnInitialize(lookupConfigFiles)

	cmd.AddGroup(&cobra.Group{
		ID:    "domains",
		Title: "Domains commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    "identity",
		Title: "Identity commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    "certs",
		Title: "Certificates commands",
	})

	cmd.AddCommand(NewCmdDomain(opts))
	cmd.AddCommand(NewCmdDSR(opts))
	cmd.AddCommand(NewCmdDNSSEC(opts))
	cmd.AddCommand(NewCmdCollaborator(opts))
	cmd.AddCommand(NewCmdEmailForward(opts))
	cmd.AddCommand(NewCmdWhoami(opts))
	cmd.AddCommand(NewCmdPush(opts))
	cmd.AddCommand(NewCmdCert(opts))
	cmd.AddCommand(NewCmdTLD(opts))
	cmd.AddCommand(NewCmdLetsEncrypt(opts))
	cmd.AddCommand(NewCmdTransfer(opts))
	cmd.AddCommand(NewCmdAccounts(opts))
	cmd.AddCommand(NewCmdConfig(opts))

	cmd.PersistentFlags().String(flagAccount, "", "Account")
	cmd.PersistentFlags().String(flagBaseURL, "", "Base URL")
	cmd.PersistentFlags().String(flagAccessToken, "", "Access token")
	cmd.PersistentFlags().Bool(flagSandbox, false, "Sandbox environment")
	cmd.PersistentFlags().StringVarP(&configFile, flagConfig, "c", "", "Configuration file")
	cmd.PersistentFlags().StringVar(&profile, flagProfile, defaultProfile, "Profile")

	cmd.MarkFlagsMutuallyExclusive(flagBaseURL, flagSandbox)

	viper.SetEnvPrefix(envPrefix)
	_ = viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}

func lookupConfigFiles() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, "dnsimple"))
		viper.AddConfigPath("/etc/dnsimple")
		viper.SetConfigType(defaultConfigFileFormat)
		viper.SetConfigName(profile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Found error: ", err.Error())
		}
	}
}
