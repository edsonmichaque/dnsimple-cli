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

package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type Account []dnsimple.Account

func (a Account) Columns() []string {
	return []string{
		"ID",
		"EMAIL",
		"PLAN IDENTIFIER",
		"CREATED AT",
		"UPDATED AT",
	}
}

func (a Account) Data() []map[string]string {
	data := make([]map[string]string, 0, len(a))

	for _, k := range a {
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

func (a Account) JSON() (io.Reader, error) {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
