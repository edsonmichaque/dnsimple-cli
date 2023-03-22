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

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/formatter"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CmdDomain(opts *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "Manage domains",
		Aliases: []string{"domains"},
	}

	cmd.AddCommand(CmdDomainCollaborator(opts))
	cmd.AddCommand(CmdDomainCreate(opts))
	cmd.AddCommand(CmdDomainDNSSec(opts))
	cmd.AddCommand(CmdDomainDSR(opts))
	cmd.AddCommand(CmdDomainDSR(opts))
	cmd.AddCommand(CmdDomainDelete(opts))
	cmd.AddCommand(CmdDomainGet(opts))
	cmd.AddCommand(CmdDomainList(opts))

	return cmd
}

func CmdDomainList(opts *Options) *cobra.Command {
	cmd := createCommand(&cobra.Command{
		Use:   "list",
		Short: "List domains",
		Args:  cobra.NoArgs,
		Example: heredoc.Doc(`
			dnsimple domain list
			dnsimple domain list --sandbox
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

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := apiClient.Domains.ListDomains(context.Background(), cfg.Account, &dnsimple.DomainListOptions{
				ListOptions: getListOptions(),
			})
			if err != nil {
				return err
			}

			output := viper.GetString("output")
			if output != formatTable && output != formatJSON && output != formatYAML {
				return errors.New("invalid output format" + output)
			}

			formattedOutput, err := formatter.Format(formatter.DomainList(*resp), &formatter.Options{
				Format: formatter.OutputFormat(output),
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

	addOutputFlag(cmd, formatTable)
	addPaginationFlags(cmd)
	addQueryFlag(cmd)

	return cmd
}

func CmdDomainDelete(opts *Options) *cobra.Command {
	cmd := createCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a domain",
		Args:  cobra.NoArgs,
		Example: heredoc.Doc(`
			dnsimple domain delete --domain example.com
			dnsimple domain delete --domain example.com --sandbox
		`),
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				panic(err)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			confirm := viper.GetBool(flagConfirm)

			if !confirm {
				if !runConfirm(viper.GetString(flagDomain)) {
					return errors.New("no confirmation")
				}
			}

			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString(flagDomain)

			if domain == "" {
				domain, err = runPromptDomainName()
				if err != nil {
					return nil
				}
			}

			_, err = opts.BuildClient(cfg.BaseURL, cfg.AccessToken).Domains.DeleteDomain(
				context.Background(),
				cfg.Account,
				domain)
			if err != nil {
				return err
			}

			cmd.Printf("✓ Deleted domain %s", domain)

			return nil
		},
	}, opts)

	addDomainFlag(cmd)
	addConfirmFlag(cmd)

	return cmd
}

func CmdDomainCreate(opts *Options) *cobra.Command {
	cmd := createCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a domain",
		Example: heredoc.Doc(`
			dnsimple domain create --domain example.com
			dnsimple domain create --domain example.com --sandbox
		`),
		Args: cobra.MaximumNArgs(1),
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

			var domain dnsimple.Domain
			if len(args) != 0 {
				err = json.Unmarshal([]byte(args[0]), &domain)
				if err != nil {
					return err
				}
			}

			if domain == (dnsimple.Domain{}) {
				domain.Name, err = runPromptDomainName()
				if err != nil {
					return nil
				}
			}

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

			resp, err := apiClient.Domains.CreateDomain(
				context.Background(),
				cfg.Account,
				domain)
			if err != nil {
				return err
			}

			cmd.Printf("%s Created domain %s\n", color.GreenString("✓"), resp.Data.Name)

			return nil
		},
	}, opts)

	return cmd
}

func CmdDomainGet(opts *Options) *cobra.Command {
	cmd := createCommand(&cobra.Command{
		Use:   "get",
		Short: "Retrieve a domain",
		Args:  cobra.NoArgs,
		Example: heredoc.Doc(`
			dnsimple domain show --domain example.com
			dnsimple domain show --domain example.com --sandbox
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

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)
			resp, err := apiClient.Domains.GetDomain(
				context.Background(),
				cfg.Account,
				viper.GetString("domain"),
			)
			if err != nil {
				return err
			}

			output := viper.GetString("output")
			if output != formatText && output != formatJSON && output != formatYAML {
				return errors.New("invalid output format")
			}

			formattedOutput, err := formatter.Format(formatter.DomainItem(*resp), &formatter.Options{
				Format: formatter.OutputFormat(output),
				// TODO: query should be only used for JSON and YAML output formats
				Query: viper.GetString("query"),
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

	addQueryFlag(cmd)
	addOutputFlag(cmd, formatText)

	return cmd
}

func addDomainFlag(cmd *cobra.Command) {
	cmd.Flags().String("domain", "", "Domain flags")
	if err := cmd.MarkFlagRequired("domain"); err != nil {
		panic(err)
	}
}

func addConfirmFlag(cmd *cobra.Command) {
	cmd.Flags().Bool("confirm", false, "Confirm")
}

func runPromptDomainName() (string, error) {
	prompt := &survey.Input{
		Message: "Domain name",
	}

	var domain string

	if err := survey.AskOne(prompt, &domain); err != nil {
		return "", err
	}

	return domain, nil
}

func addPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 0, "Page")
	cmd.Flags().Int("per-page", 0, "Per page")
}

func getListOptions() dnsimple.ListOptions {
	var opts dnsimple.ListOptions

	if page := viper.GetInt("page"); page != 0 {
		opts.Page = &page
	}

	if perPage := viper.GetInt("per-page"); perPage != 0 {
		opts.PerPage = &perPage
	}

	return opts
}

func getListOptionsP() *dnsimple.ListOptions {
	opts := getListOptions()

	return &opts
}
