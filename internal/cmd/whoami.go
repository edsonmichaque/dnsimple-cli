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

	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdWhoami(opts *internal.CmdOpt) *cobra.Command {
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

			cmd.Printf("%#q\n", resp.Data)

			return nil
		},

		GroupID: "identity",
	}

	return cmd
}
