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

	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/formatter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CmdWhoami(opts *Options) *cobra.Command {
	v := viper.New()

	cmd := createCommand(&cobra.Command{
		Use:   "whoami",
		Short: "Check identity",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.New()
			if err != nil {
				return err
			}

			resp, err := opts.BuildClient(cfg.BaseURL, cfg.AccessToken).Identity.Whoami(context.Background())
			if err != nil {
				return err
			}

			output := v.GetString(flagOutput)
			if output != "text" && output != "json" && output != "yaml" {
				return errors.New("invalid output format")
			}

			formattedOutput, err := formatter.Format(formatter.Whoami(*resp), &formatter.Options{
				Format: formatter.OutputFormat(output),
				// TODO: query should be only used for JSON and YAML output formats
				Query: v.GetString(flagQuery),
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

	addOutputFlag(cmd, formatText)
	addQueryFlag(cmd)

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	return cmd
}
