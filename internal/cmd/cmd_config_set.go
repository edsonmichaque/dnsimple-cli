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
	"github.com/edsonmichaque/dnsimple-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

var configPropValidate = map[string]func(string) (interface{}, error){
	"sandbox": func(value string) (interface{}, error) {
		return strconv.ParseBool(value)
	},
	"account": func(value string) (interface{}, error) {
		return strconv.ParseInt(value, 10, 64)
	},
	"base-url": func(value string) (interface{}, error) {
		return value, nil
	},
	"access-token": func(value string) (interface{}, error) {
		return value, nil
	},
}

func NewCmdConfigSet(opts *internal.CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Manage configurations",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, ok := configProps[args[0]]; !ok {
				return errors.New("not found")
			}

			validate := configPropValidate[args[0]]
			if validate == nil {
				return errors.New("no validator found")
			}

			value, err := validate(args[1])
			if err != nil {
				return err
			}

			viper.Set(args[0], value)

			if err := viper.WriteConfig(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
