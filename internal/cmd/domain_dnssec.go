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
	"github.com/spf13/viper"
)

func NewCmdDomainDNSSec(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dnssec",
		Short: "Manage DNSSEC status",
	}

	cmd.AddCommand(NewCmdDNSSECStatus(opts))
	cmd.AddCommand(NewCmdDNSSECDisable(opts))
	cmd.AddCommand(NewCmdDNSSECEnable(opts))

	addDomainPersistentFlag(cmd)

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		panic(err)
	}

	return cmd
}

func NewCmdDNSSECStatus(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Retrieve DNSSEC status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString("domain")

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

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
	}

	return cmd
}

func NewCmdDNSSECDisable(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable DNSSEC",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString("domain")

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

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
	}

	return cmd
}

func NewCmdDNSSECEnable(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable DNSSEC",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			cfg, err := config.New()
			if err != nil {
				return err
			}

			domain := viper.GetString("domain")

			apiClient := opts.BuildClient(cfg.BaseURL, cfg.AccessToken)

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
	}

	return cmd
}
