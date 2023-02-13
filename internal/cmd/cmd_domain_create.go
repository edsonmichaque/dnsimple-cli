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
	"fmt"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdDomainCreate(opts *internal.CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			resp, err := opts.BuildClient(cfg.BaseURL, cfg.AccessToken).Domains.CreateDomain(context.Background(), cfg.Account, dnsimple.Domain{Name: args[0]})
			if err != nil {
				return err
			}

			fmt.Printf("%#v", resp.Data)

			return nil
		},
	}

	return cmd
}
