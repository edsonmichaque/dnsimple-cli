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

	"github.com/edsonmichaque/dnsimple-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CmdDomainDNSSec(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   "dnssec",
		Short: "Manage DNSSEC status",
		Args:  cobra.NoArgs,
	}, opts)

	cmd.AddCommand(CmdDnssecStatus(opts))
	cmd.AddCommand(CmdDnssecDisable(opts))
	cmd.AddCommand(CmdDnssecEnable(opts))

	addDomainRequiredFlag(cmd)

	return cmd
}

func CmdDnssecStatus(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   "status",
		Short: "Retrieve DNSSEC status",
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

			resp, err := apiClient.Domains.GetDnssec(
				context.Background(),
				cfg.Account,
				domain,
			)
			if err != nil {
				return err
			}

			status := "disabled"
			if resp.Data.Enabled {
				status = "enabled"
			}

			cmd.Printf("✓ DNSSEC for %v is %v", domain, status)

			return nil
		},
	}, opts)

	return cmd
}

func CmdDnssecDisable(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   "disable",
		Short: "Disable DNSSEC",
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

			_, err = apiClient.Domains.DisableDnssec(
				context.Background(),
				cfg.Account,
				domain,
			)
			if err != nil {
				return err
			}

			cmd.Printf("✓ DNSSEC for %v has been disabled", domain)

			return nil
		},
	}, opts)

	return cmd
}

func CmdDnssecEnable(opts *Options) *cobra.Command {
	cmd := createCmd(&cobra.Command{
		Use:   "enable",
		Short: "Enable DNSSEC",
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

			_, err = apiClient.Domains.EnableDnssec(
				context.Background(),
				cfg.Account,
				domain,
			)
			if err != nil {
				return err
			}

			cmd.Printf("✓ DNSSEC for %v has been enabled", domain)

			return nil
		},
	}, opts)

	return cmd
}
