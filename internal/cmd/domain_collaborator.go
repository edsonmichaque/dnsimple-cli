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

	addDomainPersistentFlag(cmd)

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
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
func addFormFileFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("from-file", "f", "", "Create from file")
}

func NewCmdCollaboratorAdd(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add collaborator",
		Args:  cobra.MaximumNArgs(1),
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

			var attr dnsimple.CollaboratorAttributes

			var data []byte

			if len(args) == 0 {
				fromFile := viper.GetString("from-file")
				if fromFile == "" {
					return errors.New("data is required")
				}

				if fromFile == "-" {
					data, err = io.ReadAll(cmd.InOrStdin())
					if err != nil {
						return err
					}
				} else {
					data, err = os.ReadFile(fromFile)
					if err != nil {
						return err
					}
				}
			}

			if len(args) != 0 {
				data = []byte(args[0])
			}

			err = json.Unmarshal(data, &attr)
			if err != nil {
				return err
			}

			resp, err := opts.BuildClient(cfg.BaseURL, cfg.AccessToken).Domains.AddCollaborator(
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

	addFormFileFlag(cmd)

	if err := viper.BindPFlag("from-file", cmd.Flags().Lookup("from-file")); err != nil {
		panic(err)
	}

	return cmd
}
