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
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/printer"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdDomainDSR(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dsr",
		Short: "Manage domain signed records",
	}

	cmd.AddCommand(NewCmdDomainDSRCreate(opts))
	cmd.AddCommand(NewCmdDomainDSRList(opts))
	cmd.AddCommand(NewCmdDomainDSRGet(opts))

	addDomainRequiredFlag(cmd)

	return cmd
}

func NewCmdDomainDSRCreate(opts *internal.CommandOptions) *cobra.Command {
	v := viper.New()

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a delegation signer record",
		Example: heredoc.Doc(`
			dnsimple domain dsr create --domain example.com
			dnsimple domain dsr create --domain example.com --sandbox
		`),
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				panic(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			var (
				domain   = viper.GetString("domain")
				fromFile = v.GetString("from-file")
			)

			var rawBody []byte

			if len(args) != 0 {
				rawBody = []byte(args[0])
			}

			if len(args) == 0 {
				if fromFile == "-" {
					rawBody, err = io.ReadAll(cmd.InOrStdin())
					if err != nil {
						return err
					}
				}

				if fromFile != "" && fromFile != "-" {
					rawBody, err = os.ReadFile(fromFile)
					if err != nil {
						return err
					}
				}
			}

			if len(rawBody) == 0 {
				return errors.New("body is required")
			}

			var attr dnsimple.DelegationSignerRecord

			err = json.Unmarshal(rawBody, &attr)
			if err != nil {
				return err
			}

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := apiClient.Domains.CreateDelegationSignerRecord(
				context.Background(),
				cfg.Account,
				domain,
				attr,
			)
			if err != nil {
				return err
			}

			cmd.Printf("%s Added delegation signer record %v\n", color.GreenString("✓"), resp.Data.ID)

			return nil
		},
	}

	addFromFileFlag(cmd)

	return cmd
}

func NewCmdDomainDSRList(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List delegation signer records",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				panic(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString("domain")

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := apiClient.Domains.ListDelegationSignerRecords(context.Background(), cfg.Account, domain, getListOptionsP())
			if err != nil {
				return err
			}

			output := viper.GetString("output")
			if output != formatTable && output != formatJSON && output != formatYAML {
				return errors.New("invalid output format")
			}

			reader, err := printer.Print(printer.DSRList(*resp), &printer.Options{
				Format: printer.Format(output),
				// TODO: query should be only used for JSON and YAML output formats
				Query: viper.GetString("query"),
			})
			if err != nil {
				return err
			}

			if _, err := io.Copy(os.Stdout, reader); err != nil {
				return err
			}

			return nil
		},
	}

	addPaginationFlags(cmd)
	addQueryFlag(cmd)
	addOutputFlag(cmd, "table")

	return cmd
}

func NewCmdDomainDSRGet(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Retrieve a delegation signer record",
		Example: heredoc.Doc(`
			$ dnsimple dsr get --domain example.com --record-id 1
			$ dnsimple dsr get --domain example.com --record-id 1 --sandbox
		`),
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				panic(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)
			resp, err := apiClient.Domains.GetDelegationSignerRecord(
				context.Background(),
				cfg.Account,
				viper.GetString("domain"),
				viper.GetInt64("record-id"),
			)
			if err != nil {
				return err
			}

			output := viper.GetString("output")
			if output != formatText && output != formatJSON && output != formatYAML {
				return errors.New("invalid output format")
			}

			printData, err := printer.Print(printer.DSRItem(*resp), &printer.Options{
				Format: printer.Format(output),
				// TODO: query should be only used for JSON and YAML output formats
				Query: viper.GetString("query"),
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

	addRecordIDFlag(cmd)
	addQueryFlag(cmd)
	addOutputFlag(cmd, "text")

	return cmd
}

func addRecordIDFlag(cmd *cobra.Command) {
	cmd.Flags().String("record-id", "", "Record id")
	if err := cmd.MarkFlagRequired("record-id"); err != nil {
		panic(err)
	}
}
