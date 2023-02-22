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

package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type DSRItem dnsimple.DelegationSignerRecordResponse

func (d DSRItem) FormatText(opts *Options) (io.Reader, error) {
	keys := []string{
		"id",
		"domain_id",
		"algorithm",
		"digest",
		"digest_type",
		"keytag",
		"public_key",
		"created_at",
		"updated_at",
	}

	const txtLen = 8

	values := map[string]interface{}{
		"id":          d.Data.ID,
		"domain_id":   d.Data.DomainID,
		"algorithm":   d.Data.Algorithm,
		"digest":      truncate(d.Data.Digest, txtLen),
		"digest_type": d.Data.DigestType,
		"keytag":      d.Data.Keytag,
		"public_key":  truncate(d.Data.PublicKey, txtLen),
		"created_at":  d.Data.CreatedAt,
		"updated_at":  d.Data.UpdatedAt,
	}

	titles := map[string]string{
		"id":          "ID",
		"domain_id":   "Domain ID",
		"algorithm":   "Algorithm",
		"digest":      "Digest",
		"digest_type": "Digest type",
		"keytag":      "Keytag",
		"public_key":  "Public key",
		"created_at":  "Created at",
		"updated_at":  "Updated at",
	}

	buf := new(bytes.Buffer)
	for _, v := range keys {
		buf.WriteString(fmt.Sprintf("%-20s%v\n", titles[v]+":", values[v]))
	}

	return buf, nil
}

func (d DSRItem) FormatJSON(opts *Options) (io.Reader, error) {
	return formatJSON(d, opts)
}

func (d DSRItem) FormatYAML(opts *Options) (io.Reader, error) {
	return formatYAML(d, opts)
}

func (d DSRItem) formatJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(d.Data, "", "  ")
}
