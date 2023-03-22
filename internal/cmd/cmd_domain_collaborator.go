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
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/formatter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CmdDomainCollaborator(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "collaborator",
		Args: cobra.NoArgs,
		Example: heredoc.Doc(`
			dnsimple collaborator
		`),
		Short: "Manage collaborators",
	}

	cmd.AddCommand(CmdCollaboratorAdd(opts))
	cmd.AddCommand(CmdCollaboratorList(opts))
	cmd.AddCommand(CmdCollaboratorRemove(opts))

	addDomainRequiredFlag(cmd)

	return cmd
}

func CmdCollaboratorAdd(opts *Options) *cobra.Command {
	cmd := createCommand(&cobra.Command{
		Use:   "add --domain [DOMAIN]",
		Short: "Add collaborator",
		Args:  cobra.MaximumNArgs(1),
		Example: heredoc.Doc(`
			dnsimple collaborator add --domain example.com '{"email":"john.doe@example.com"}'
			dnsimple collaborator add --domain example.com --from-file email.json
			dnsimple collaborator add --domain example.com --from-file=-
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
				domain   = viper.GetString(flagDomain)
				fromFile = viper.GetString(flagFromFile)
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
	}, opts)

	addFromFileFlag(cmd)

	return cmd
}

func CmdCollaboratorRemove(opts *Options) *cobra.Command {
	cmd := createCommand(&cobra.Command{
		Use:   "remove --domain [DOMAIN]",
		Short: "Remove collaborator",
		Args:  cobra.NoArgs,
		Example: heredoc.Doc(`
			dnsimple collaborator
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

			domain := viper.GetString(flagDomain)

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			_, err = apiClient.Domains.RemoveCollaborator(
				context.Background(),
				cfg.Account,
				domain,
				viper.GetInt64(flagCollaboratorID),
			)
			if err != nil {
				return err
			}

			cmd.Printf("✓ Deleted collaborator %v", viper.GetInt64(flagCollaboratorID))

			return nil
		},
	}, opts)

	addDomainRequiredFlag(cmd)
	addCollaboratorIDFlag(cmd)

	return cmd
}

func addCollaboratorIDFlag(cmd *cobra.Command) {
	cmd.Flags().Int64("collaborator-id", 0, "Collaborator id")
	if err := cmd.MarkFlagRequired("collaborator-id"); err != nil {
		panic(err)
	}
}

func applyOpts(cmd *cobra.Command, opts *Options) {
	if opts.Stdout != nil {
		cmd.SetOut(opts.Stdout)
	}

	if opts.Stdin != nil {
		cmd.SetIn(opts.Stdin)
	}

	if opts.Stderr != nil {
		cmd.SetErr(opts.Stderr)
	}
}

func CmdCollaboratorList(opts *Options) *cobra.Command {
	cmd := createCommand(&cobra.Command{
		Use: "list --domain [DOMAIN]",
		Example: heredoc.Doc(`
			dnsimple collaborator list --domain example.com
			dnsimple collaborator list --domain example.com
		`),
		Args:  cobra.NoArgs,
		Short: "List collaborators",
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

			domain := viper.GetString(flagDomain)

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := apiClient.Domains.ListCollaborators(context.Background(), cfg.Account, domain, nil)
			if err != nil {
				return err
			}

			output := viper.GetString(flagOutput)
			if output != formatTable && output != formatJSON && output != formatYAML {
				return errors.New("invalid output format")
			}

			reader, err := formatter.Format(formatter.CollaboratorList(*resp), &formatter.Options{
				Format: formatter.OutputFormat(output),
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

	addQueryFlag(cmd)
	addOutputFlag(cmd, formatTable)
	addPaginationFlags(cmd)

	return cmd
}

func addDomainRequiredFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(flagDomain, "", "Domain name")
	if err := cmd.MarkPersistentFlagRequired(flagDomain); err != nil {
		panic(err)
	}
}

func addFromFileFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(flagFromFile, "f", "", "Create from file")
}
