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
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 0, "Page")
	cmd.Flags().Int("per-page", 0, "Per page")
}

func addQueryFlags(cmd *cobra.Command) {
	cmd.Flags().String("query", "", "Query")
}

func addFormatFlags(cmd *cobra.Command, format string) {
	cmd.Flags().String("format", format, "Format")
}

func checkQueryFlag(v *viper.Viper) error {
	if v.GetString("format") != "json" && v.GetString("query") != "" {
		return errors.New("query flag can only be used with json format")
	}

	return nil
}
