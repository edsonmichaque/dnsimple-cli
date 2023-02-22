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

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdDomainCollaborator(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use: "collaborator",
		Example: heredoc.Doc(`
			dnsimple collaborator
		`),
		Short: "Manage collaborators",
	}

	cmd.AddCommand(NewCmdCollaboratorAdd(opts))
	cmd.AddCommand(NewCmdCollaboratorList(opts))
	cmd.AddCommand(NewCmdCollaboratorRemove(opts))

	addDomainPersistentFlag(cmd)

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		panic(err)
	}

	return cmd
}

func NewCmdCollaboratorAdd(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add --domain [DOMAIN]",
		Short: "Add collaborator",
		Args:  cobra.MaximumNArgs(1),
		Example: heredoc.Doc(`
			dnsimple collaborator add --domain example.com '{"email":"john.doe@example.com"}'
			dnsimple collaborator add --domain example.com --from-file email.json
			dnsimple collaborator add --domain example.com --from-file=-
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			var (
				domain   = viper.GetString("domain")
				fromFile = viper.GetString("from-file")
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

			var attr dnsimple.CollaboratorAttributes

			err = json.Unmarshal(rawBody, &attr)
			if err != nil {
				return err
			}

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := apiClient.Domains.AddCollaborator(
				context.Background(),
				cfg.Account,
				domain,
				attr,
			)
			if err != nil {
				return err
			}

			cmd.Printf("%s Added collaborator %s\n", color.GreenString("✓"), resp.Data.UserEmail)

			return nil
		},
	}

	addFromFileFlag(cmd)

	if err := viper.BindPFlag("from-file", cmd.Flags().Lookup("from-file")); err != nil {
		panic(err)
	}

	return cmd
}

func NewCmdCollaboratorRemove(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove --domain [DOMAIN]",
		Short: "Remove collaborator",
		Args:  cobra.NoArgs,
		Example: heredoc.Doc(`
			dnsimple collaborator
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString("domain")

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			_, err = apiClient.Domains.RemoveCollaborator(
				context.Background(),
				cfg.Account,
				domain,
				viper.GetInt64("collaborator-id"),
			)
			if err != nil {
				return err
			}

			cmd.Printf("✓ Deleted collaborator %v", viper.GetInt64("collaborator-id"))

			return nil
		},
	}

	addCollaboratorIDFlag(cmd)

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	return cmd
}

func addCollaboratorIDFlag(cmd *cobra.Command) {
	cmd.Flags().Int64("collaborator-id", 0, "Collaborator id")
	if err := cmd.MarkFlagRequired("collaborator-id"); err != nil {
		panic(err)
	}
}

func NewCmdCollaboratorList(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use: "list --domain [DOMAIN]",
		Example: heredoc.Doc(`
			dnsimple collaborator list --domain example.com
			dnsimple collaborator list --domain example.com
		`),
		Short: "List collaborators",
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString("domain")

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := apiClient.Domains.ListCollaborators(context.Background(), cfg.Account, domain, nil)
			if err != nil {
				return err
			}

			output := viper.GetString("output")
			if output != formatTable && output != formatJSON && output != formatYAML {
				return errors.New("invalid output format")
			}

			reader, err := printer.Print(printer.CollaboratorList(*resp), &printer.Options{
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

	addQueryFlag(cmd)
	addOutputFlag(cmd, "table")
	addPaginationFlags(cmd)

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	return cmd
}

func addDomainPersistentFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String("domain", "", "Domain name")
	if err := cmd.MarkPersistentFlagRequired("domain"); err != nil {
		panic(err)
	}
}
func addFromFileFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("from-file", "f", "", "Create from file")
}
