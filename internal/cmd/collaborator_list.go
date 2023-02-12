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
	"io"
	"os"

	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/edsonmichaque/dnsimple-cli/internal/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdCollaboratorList(opts *internal.CmdOpt) *cobra.Command {
	viper := viper.New()

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List collaborators",
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString("domain")

			resp, err := opts.BuildClient(cfg.BaseURL, cfg.AccessToken).Domains.ListCollaborators(context.Background(), cfg.Account, domain, nil)
			if err != nil {
				return err
			}

			reader, err := printer.Print(printer.Collaborators(resp.Data), printer.Options{
				Format: printer.Format(viper.GetString("format")),
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

	cmd.Flags().String("domain", "", "Domain")
	cmd.Flags().String("format", "table", "Domain")

	_ = cmd.MarkFlagRequired("domain")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}
