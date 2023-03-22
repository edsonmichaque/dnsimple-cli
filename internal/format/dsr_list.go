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

type DSRList dnsimple.DelegationSignerRecordsResponse

func (a DSRList) FormatJSON(opts *Options) (io.Reader, error) {
	return formatJSON(a, opts)
}

func (a DSRList) FormatYAML(opts *Options) (io.Reader, error) {
	return formatYAML(a, opts)
}

func (a DSRList) FormatTable(_ *Options) (io.Reader, error) {
	return formatTable(a)
}

func (a DSRList) formatJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a DSRList) formatHeader() []string {
	return []string{
		"ID",
		"DOMAIN ID",
		"ALGORITHM",
		"DIGEST",
		"DIGEST TYPE",
		"KEYTAG",
		"PUBLIC KEY",
		"CREATED AT",
		"UPDATED AT",
	}
}

func (a DSRList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(a.Data))

	dsr := a.Data

	const txtLen = 10

	for i := range dsr {
		data = append(data, map[string]string{
			"ID":          fmt.Sprintf("%d", dsr[i].ID),
			"DOMAIN ID":   fmt.Sprintf("%d", dsr[i].DomainID),
			"ALGORITHM":   dsr[i].Algorithm,
			"DIGEST":      truncate(dsr[i].Digest, txtLen),
			"DIGEST TYPE": dsr[i].DigestType,
			"KEYTAG":      dsr[i].Keytag,
			"PUBLIC KEY":  truncate(dsr[i].PublicKey, txtLen),
			"CREATED AT":  dsr[i].CreatedAt,
			"UPDATED AT":  dsr[i].UpdatedAt,
		})
	}

	return data
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}

	return s[:length] + "..."
}
