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
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdDomain(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "Manage domains",
		Aliases: []string{"domains"},
	}

	cmd.AddCommand(NewCmdDomainList(opts))

	return cmd
}

func NewCmdDomainList(opts *internal.CommandOptions) *cobra.Command {
	v := viper.New()

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List domains",
		Example: heredoc.Doc(`
			dnsimple domain list
			dnsimple domain list --sandbox
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			client := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := client.Domains.ListDomains(context.Background(), cfg.Account, &dnsimple.DomainListOptions{
				ListOptions: getListOptions(v),
			})
			if err != nil {
				return err
			}

			output := v.GetString("output")
			if output != "table" && output != "json" && output != "yaml" {
				return errors.New("invalid output format")
			}

			printData, err := printer.Print(printer.DomainList(*resp), &printer.Options{
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

	addPaginationFlags(cmd)
	addQueryFlag(cmd)
	addOutputFlag(cmd, "table")

	if err := v.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	return cmd
}

func addPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 0, "Page")
	cmd.Flags().Int("per-page", 0, "Per page")
}

func getListOptions(v *viper.Viper) dnsimple.ListOptions {
	var opts dnsimple.ListOptions

	if page := v.GetInt("page"); page != 0 {
		opts.Page = &page
	}

	if perPage := v.GetInt("per-page"); perPage != 0 {
		opts.PerPage = &perPage
	}

	return opts
}
