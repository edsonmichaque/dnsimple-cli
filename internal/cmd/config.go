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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.NewWithValidation(false)
			if err != nil {
				return err
			}

			home, err := os.UserConfigDir()
			if err != nil {
				return err
			}

			profile := viper.GetString("profile")

			cmd.Println(fmt.Sprintf("Configuring profile '%s'", profile))
			cfg, ext, err := promptConfig(cfg)
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
			if err := v.WriteConfigAs(cfgPath); err != nil {
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
		Use:   "get",
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

const (
	envDev     = "DEV"
	envSandbox = "SANDBOX"
	envProd    = "PROD"
)

func promptConfig(c *config.Config) (*config.Config, string, error) {
	baseURL := prodBaseURL

	accountID, err := promptAccountID(c.Account)
	if err != nil {
		return nil, "", err
	}

	accessToken, err := promptAccessToken(c.AccessToken)
	if err != nil {
		return nil, "", err
	}

	env, err := promptEnvironment(envProd)
	if err != nil {
		return nil, "", err
	}

	if env == envSandbox {
		baseURL = sandboxBaseURL
	}

	if env == envDev {
		baseURL, err = promptBaseURL(prodBaseURL)
		if err != nil {
			return nil, "", err
		}
	}

	fileFormat, err := promptFileFormat(configFormatJSON)
	if err != nil {
		return nil, "", err
	}

	confirmation, err := promptConfirmation("Do you want to save?", true)
	if err != nil {
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

func promptAccessToken(value string) (string, error) {
	prompt := &survey.Input{
		Message: "Access Token",
		Default: value,
	}

	var token string
	if err := survey.AskOne(prompt, &token); err != nil {
		return "", err
	}

	return token, nil
}

func promptAccountID(value string) (string, error) {
	prompt := &survey.Input{
		Message: "Account ID",
		Default: value,
	}

	var accountID string
	if err := survey.AskOne(prompt, &accountID); err != nil {
		return "", err
	}

	return accountID, nil
}

func promptEnvironment(value string) (string, error) {
	prompt := &survey.Select{
		Message: "Environment",
		Options: []string{envProd, envSandbox, envDev},
		Default: value,
	}

	var env string
	if err := survey.AskOne(prompt, &env); err != nil {
		return "", err
	}

	return env, nil
}

func promptBaseURL(value string) (string, error) {
	prompt := &survey.Input{
		Message: "Base URL",
		Default: value,
	}

	var baseURL string
	if err := survey.AskOne(prompt, &baseURL); err != nil {
		return "", err
	}

	return baseURL, nil
}

func promptFileFormat(value string) (string, error) {
	prompt := &survey.Select{
		Message: "File format",
		Options: []string{configFormatJSON, configFormatYAML, configFormatTOML},
		Default: value,
	}

	var fileFormat string
	if err := survey.AskOne(prompt, &fileFormat); err != nil {
		return "", err
	}

	return fileFormat, nil
}

func promptConfirmation(msg string, value bool) (bool, error) {
	var confirmation bool

	prompt := &survey.Confirm{
		Message: msg,
		Default: value,
	}

	if err := survey.AskOne(prompt, &confirmation); err != nil {
		return false, err
	}

	return confirmation, nil
}
