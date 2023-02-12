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
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FormatJSON = "json"
	FormatYAML = "yaml"
	FormatTOML = "toml"
	FormatYML  = "yml"
)

func NewCmdConfig(opts *internal.CmdOpt) *cobra.Command {
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
				FormatJSON,
				FormatYAML,
				FormatTOML,
				FormatYML,
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

	return cmd
}

func runConfigPrompt() (*config.Config, string, error) {
	promptAccountID := &survey.Input{
		Message: "Account ID",
	}

	promptAccessToken := &survey.Password{
		Message: "Access Token",
	}

	promptEnv := &survey.Select{
		Message: "Environment",
		Options: []string{"PROD", "SANDBOX", "DEV"},
		Default: "PROD",
	}

	var (
		accountID   string
		accessToken string
		env         string
		baseURL     = "https://api.dnsimple.com"
	)

	if err := survey.AskOne(promptAccountID, &accountID); err != nil {
		return nil, "", err
	}

	if err := survey.AskOne(promptAccessToken, &accessToken); err != nil {
		return nil, "", err
	}

	if err := survey.AskOne(promptEnv, &env); err != nil {
		return nil, "", err
	}

	if env == "SANDBOX" {
		baseURL = "https://api.sandbox.dnsimple.com"
	}

	if env == "DEV" {
		promptBaseURL := &survey.Input{
			Message: "Base URL",
		}

		if err := survey.AskOne(promptBaseURL, &baseURL, survey.WithValidator(validateURL)); err != nil {
			return nil, "", err
		}
	}

	promptFileFormat := &survey.Select{
		Message: "File format",
		Options: []string{"JSON", "YAML", "TOML"},
		Default: "JSON",
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

	if env == "DEV" {
		cfg.BaseURL = baseURL
	}

	if env == "SANDBOX" {
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
