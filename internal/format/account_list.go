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

package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type Filter struct {
	Fields []string
}

type AccountList dnsimple.AccountsResponse

func (a AccountList) FormatJSON(opts *Options) (io.Reader, error) {
	return formatJSON(a, opts)
}

func (a AccountList) FormatYAML(opts *Options) (io.Reader, error) {
	return formatYAML(a, opts)
}

func (a AccountList) FormatTable(_ *Options) (io.Reader, error) {
	return formatTable(a)
}

func (a AccountList) formatJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a AccountList) formatHeader() []string {
	return []string{
		"ID",
		"EMAIL",
		"PLAN IDENTIFIER",
		"CREATED AT",
		"UPDATED AT",
	}
}

func (a AccountList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(a.Data))

	for _, k := range a.Data {
		data = append(data, map[string]string{
			"ID":              fmt.Sprintf("%d", k.ID),
			"EMAIL":           k.Email,
			"PLAN IDENTIFIER": k.PlanIdentifier,
			"CREATED AT":      k.CreatedAt,
			"UPDATED AT":      k.UpdatedAt,
		})
	}

	return data
}
