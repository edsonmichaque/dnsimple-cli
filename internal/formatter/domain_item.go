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

type DomainItem dnsimple.DomainResponse

func (d DomainItem) FormatText(opts *Options) (io.Reader, error) {
	keys := []string{
		"id",
		"account_id",
		"registrant_id",
		"name",
		"unicode_name",
		"token",
		"state",
		"auto_renew",
		"private_whois",
		"expires_at",
		"created_at",
		"updated_at",
	}

	values := map[string]interface{}{
		"id":            d.Data.ID,
		"account_id":    d.Data.AccountID,
		"registrant_id": d.Data.RegistrantID,
		"name":          d.Data.Name,
		"unicode_name":  d.Data.UnicodeName,
		"token":         d.Data.Token,
		"state":         d.Data.State,
		"auto_renew":    d.Data.AutoRenew,
		"private_whois": d.Data.PrivateWhois,
		"expires_at":    d.Data.ExpiresAt,
		"created_at":    d.Data.CreatedAt,
		"updated_at":    d.Data.UpdatedAt,
	}

	titles := map[string]string{
		"id":            "ID",
		"account_id":    "Account ID",
		"registrant_id": "Registrant ID",
		"name":          "Name",
		"unicode_name":  "Unicode name",
		"token":         "Token",
		"state":         "State",
		"auto_renew":    "Auto renew",
		"private_whois": "Private whois",
		"expires_at":    "Expires at",
		"created_at":    "Created at",
		"updated_at":    "Updated at",
	}

	buf := new(bytes.Buffer)
	for _, v := range keys {
		buf.WriteString(fmt.Sprintf("%-20s%v\n", titles[v]+":", values[v]))
	}

	return buf, nil
}

func (d DomainItem) FormatJSON(opts *Options) (io.Reader, error) {
	return formatJSON(d, opts)
}

func (d DomainItem) FormatYAML(opts *Options) (io.Reader, error) {
	return formatYAML(d, opts)
}

func (d DomainItem) formatJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(d.Data, "", "  ")
}
