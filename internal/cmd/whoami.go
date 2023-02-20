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

	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdWhoami(opts *internal.CmdOpt) *cobra.Command {
	v := viper.New()

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Check identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			resp, err := opts.BuildClient(cfg.BaseURL, cfg.AccessToken).Identity.Whoami(context.Background())
			if err != nil {
				return err
			}

			output := v.GetString("output")
			if output != "text" && output != "json" && output != "yaml" {
				return errors.New("invalid output format")
			}

			printData, err := printer.Print(printer.Whoami(*resp), &printer.Options{
				Format: printer.Format(output),
				// TODO: query should be only used for JSON and YAML output formats
				Query: v.GetString("query"),
			})
			if err != nil {
				return err
			}

			if _, err := io.Copy(cmd.OutOrStdout(), printData); err != nil {
				return err
			}

			return nil
		},
	}

	addOutputFlags(cmd, "text")
	addQueryFlags(cmd)

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	return cmd
}
