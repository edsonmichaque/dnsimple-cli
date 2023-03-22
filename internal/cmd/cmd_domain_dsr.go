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
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/format"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	actionList   = "list"
	actionCreate = "create"
	actionDelete = "delete"
)

func CmdDomainDSR(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   "dsr",
		Short: "Manage domain signed records",
		Args:  cobra.NoArgs,
	}, opts)

	cmd.AddCommand(CmdDomainDSRCreate(opts))
	cmd.AddCommand(CmdDomainDSRList(opts))
	cmd.AddCommand(CmdDomainDSRGet(opts))

	addDomainRequiredFlag(cmd)

	return cmd
}

func CmdDomainDSRCreate(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   actionCreate,
		Short: "Create a delegation signer record",
		Args:  cobra.MaximumNArgs(1),
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
			cfg, err := config.New()
			if err != nil {
				return err
			}

			var (
				domain   = viper.GetString(configDomain)
				fromFile = viper.GetString(optionFromFile)
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

			apiClient := opts.createClient(cfg.BaseURL, cfg.AccessToken)

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
	}, opts)

	addFromFileFlag(cmd)

	return cmd
}

func CmdDomainDSRList(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   actionList,
		Short: "List delegation signer records",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				panic(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString(configDomain)

			apiClient := opts.createClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := apiClient.Domains.ListDelegationSignerRecords(context.Background(), cfg.Account, domain, getListOptionsP())
			if err != nil {
				return err
			}

			output := viper.GetString(flagOutput)
			if output != formatTable && output != formatJSON && output != formatYAML {
				return errors.New("invalid output format")
			}

			reader, err := format.Format(format.DSRList(*resp), &format.Options{
				Format: format.OutputFormat(output),
				// TODO: query should be only used for JSON and YAML output formats
				Query: viper.GetString(flagQuery),
			})
			if err != nil {
				return err
			}

			if _, err := io.Copy(os.Stdout, reader); err != nil {
				return err
			}

			return nil
		},
	}, opts)

	addPaginationFlags(cmd)
	addQueryFlag(cmd)
	addOutputFlag(cmd, formatTable)

	return cmd
}

func CmdDomainDSRGet(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   "get",
		Short: "Retrieve a delegation signer record",
		Args:  cobra.NoArgs,
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
			cfg, err := config.New()
			if err != nil {
				return err
			}

			apiClient := opts.createClient(cfg.BaseURL, cfg.AccessToken)
			resp, err := apiClient.Domains.GetDelegationSignerRecord(
				context.Background(),
				cfg.Account,
				viper.GetString(configDomain),
				viper.GetInt64(flagRecordID),
			)
			if err != nil {
				return err
			}

			output := viper.GetString(flagOutput)
			if output != formatText && output != formatJSON && output != formatYAML {
				return errors.New("invalid output format")
			}

			formattedOutput, err := format.Format(format.DSRItem(*resp), &format.Options{
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

	addRecordIDFlag(cmd)
	addQueryFlag(cmd)
	addOutputFlag(cmd, formatText)

	return cmd
}

func addRecordIDFlag(cmd *cobra.Command) {
	cmd.Flags().String(flagRecordID, "", "Record id")
	if err := cmd.MarkFlagRequired("record-id"); err != nil {
		panic(err)
	}
}
