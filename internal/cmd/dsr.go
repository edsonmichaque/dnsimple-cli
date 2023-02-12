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
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/spf13/cobra"
)

func NewCmdDSR(opts *internal.CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dsr",
		Short:   "Manage domain signed records",
		GroupID: "domains",
	}

	cmd.AddCommand(NewCmdDSRCreate(opts))
	cmd.AddCommand(NewCmdDSRGet(opts))
	cmd.AddCommand(NewCmdDSRList(opts))
	cmd.AddCommand(NewCmdDSRDelete(opts))

	return cmd
}
