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
	"runtime"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/edsonmichaque/dnsimple-cli/internal/build"
	"github.com/spf13/cobra"
)

func NewCmdVersion(opts *internal.CommandOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Check identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SetupIO(cmd, opts)

			versionTemplate := heredoc.Doc(`
				DNSimple CLI version:       %v
				DNSimple Go Client version: %v
				DNSimple API endpoint:      %v
				DNSimple API version:       %v
				OS/Arch (client):           %v/%v
			`)

			cmd.Printf(versionTemplate, build.Version, dnsimple.Version, prodBaseURL, "v2", runtime.GOOS, runtime.GOARCH)

			return nil
		},
	}

	return cmd
}
