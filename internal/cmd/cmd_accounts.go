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
	"context"
	"errors"
	"io"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CmdAccounts(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   "accounts",
		Short: "List accounts",
		Example: heredoc.Doc(`
			dnsimple acounts
			dnsimple accounts --output=json
			dnsimple accounts --output=yaml
			dnsimple accounts --output=json --query="[].id"
		`),
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				panic(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.New()
			if err != nil {
				// TODO: pretty print error
				return err
			}

			resp, err := opts.BuildClient(cfg.BaseURL, cfg.AccessToken).Accounts.ListAccounts(context.Background(), nil)
			if err != nil {
				// TODO: pretty print error returned from the API client
				return err
			}

			output := viper.GetString(flagOutput)
			if output != "table" && output != "json" && output != "yaml" {
				return errors.New("invalid output format")
			}

			formattedOutput, err := format.Format(format.AccountList(*resp), &format.Options{
				Format: format.OutputFormat(output),
				// TODO: query should be only used for JSON and YAML output formats
				Query: viper.GetString(flagQuery),
			})
			if err != nil {
				return err
			}

			if _, err := io.Copy(cmd.OutOrStdout(), formattedOutput); err != nil {
				return err
			}

			return nil
		},
	}, opts)

	addOutputFlag(cmd, "table")
	addQueryFlag(cmd)

	return cmd
}

func addOutputFlag(cmd *cobra.Command, f string) {
	cmd.Flags().StringP(flagOutput, "o", f, "Output format")
}

func addQueryFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(flagQuery, "q", "", "Query")
}
