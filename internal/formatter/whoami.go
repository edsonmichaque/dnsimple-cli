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

type Whoami dnsimple.WhoamiResponse

func (w Whoami) FormatText(opts *Options) (io.Reader, error) {
	var data [][2]interface{}

	if user := w.Data.User; user != nil {
		data = [][2]interface{}{
			{"User", ""},
			{"  ID", user.ID},
			{"  Email", user.Email},
		}
	}

	if account := w.Data.Account; account != nil {
		data = [][2]interface{}{
			{"Accout", ""},
			{"  ID", account.ID},
			{"  Email", account.Email},
			{"  Plan identifier", account.PlanIdentifier},
			{"  Created at", account.CreatedAt},
			{"  Updated at", account.UpdatedAt},
		}
	}

	buf := new(bytes.Buffer)

	for _, i := range data {
		buf.WriteString(fmt.Sprintf("%-18s %v\n", i[0].(string)+":", i[1]))
	}

	return buf, nil
}

func (w Whoami) FormatJSON(opts *Options) (io.Reader, error) {
	return formatJSON(w, opts)
}

func (w Whoami) FormatYAML(opts *Options) (io.Reader, error) {
	return formatYAML(w, opts)
}

func (w Whoami) formatJSON(opts *Options) ([]byte, error) {
	return json.MarshalIndent(w.Data, "", "  ")
}
