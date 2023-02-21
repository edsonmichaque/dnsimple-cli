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
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFormatJSON = "json"
	configFormatYAML = "yaml"
	configFormatTOML = "toml"
	configFormatYML  = "yml"
)

var (
	prodBaseURL    = "https://api.dnsimple.com"
	sandboxBaseURL = "https://api.sandbox.dnsimple.com"

	configProps = map[string]struct{}{
		"account":      {},
		"base-url":     {},
		"access-token": {},
		"sandbox":      {},
	}

	configPropValidate = map[string]func(string) (interface{}, error){
		"sandbox": func(value string) (interface{}, error) {
			return strconv.ParseBool(value)
		},
		"account": func(value string) (interface{}, error) {
			return strconv.ParseInt(value, 10, 64)
		},
		"base-url": func(value string) (interface{}, error) {
			return value, nil
		},
		"access-token": func(value string) (interface{}, error) {
			return value, nil
		},
	}
)

func NewCmdConfig(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := os.UserConfigDir()
			if err != nil {
				return err
			}

			profile := viper.GetString("profile")

			for _, ext := range []string{
				configFormatJSON,
				configFormatYAML,
				configFormatTOML,
				configFormatYML,
			} {
				fname := fmt.Sprintf("%s.%s", profile, strings.ToLower(ext))
				if stat, _ := os.Stat(filepath.Join(home, "dnsimple", fname)); stat != nil {
					return fmt.Errorf("there is already a file for profile %s", profile)
				}
			}

			cfg, ext, err := runConfigPrompt()
			if err != nil {
				return err
			}

			v := viper.New()
			v.Set(flagAccount, cfg.Account)
			v.Set(flagAccessToken, cfg.AccessToken)
			if cfg.BaseURL != "" {
				v.Set(flagBaseURL, cfg.BaseURL)
			}

			if cfg.Sandbox {
				v.Set(flagSandbox, cfg.Sandbox)
			}

			cfgPath := filepath.Join(home, "dnsimple", fmt.Sprintf("%s.%s", profile, strings.ToLower(ext)))
			if err := v.SafeWriteConfigAs(cfgPath); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(NewCmdConfigGet(opts))
	cmd.AddCommand(NewCmdConfigSet(opts))

	return cmd
}

func NewCmdConfigGet(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Manage configurations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, ok := configProps[args[0]]; !ok {
				return errors.New("not found")
			}

			cmd.Println(viper.GetString(args[0]))

			return nil
		},
	}

	return cmd
}

func NewCmdConfigSet(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Manage configurations",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, ok := configProps[args[0]]; !ok {
				return errors.New("not found")
			}

			validate := configPropValidate[args[0]]
			if validate == nil {
				return errors.New("no validator found")
			}

			value, err := validate(args[1])
			if err != nil {
				return err
			}

			viper.Set(args[0], value)

			if err := viper.WriteConfig(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func runConfigPrompt() (*config.Config, string, error) {
	const (
		envDev     = "DEV"
		envSandbox = "SANDBOX"
		envProd    = "PROD"
	)

	promptAccountID := &survey.Input{
		Message: "Account ID",
	}

	promptAccessToken := &survey.Password{
		Message: "Access Token",
	}

	promptEnv := &survey.Select{
		Message: "Environment",
		Options: []string{envProd, envSandbox, envDev},
		Default: envProd,
	}

	var (
		accountID   string
		accessToken string
		env         string
	)

	baseURL := prodBaseURL

	if err := survey.AskOne(promptAccountID, &accountID); err != nil {
		return nil, "", err
	}

	if err := survey.AskOne(promptAccessToken, &accessToken); err != nil {
		return nil, "", err
	}

	if err := survey.AskOne(promptEnv, &env); err != nil {
		return nil, "", err
	}

	if env == envSandbox {
		baseURL = sandboxBaseURL
	}

	if env == envDev {
		promptBaseURL := &survey.Input{
			Message: "Base URL",
		}

		if err := survey.AskOne(promptBaseURL, &baseURL, survey.WithValidator(validateURL)); err != nil {
			return nil, "", err
		}
	}

	promptFileFormat := &survey.Select{
		Message: "File format",
		Options: []string{configFormatJSON, configFormatYAML, configFormatTOML},
		Default: configFormatJSON,
	}

	var fileFormat string
	if err := survey.AskOne(promptFileFormat, &fileFormat); err != nil {
		return nil, "", err
	}

	promptConfirm := &survey.Confirm{
		Message: "Do you want to save?",
	}

	var confirmation bool
	if err := survey.AskOne(promptConfirm, &confirmation); err != nil {
		return nil, "", err
	}

	if !confirmation {
		return nil, "", errors.New("did not confirm")
	}

	cfg := config.Config{
		Account:     accountID,
		AccessToken: accessToken,
	}

	if env == envDev {
		cfg.BaseURL = baseURL
	}

	if env == envSandbox {
		cfg.Sandbox = true
	}

	return &cfg, fileFormat, nil
}

func validateURL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return errors.New("no string")
	}

	if _, err := url.Parse(str); err != nil {
		return err
	}

	fmt.Println("No error")

	return nil
}
