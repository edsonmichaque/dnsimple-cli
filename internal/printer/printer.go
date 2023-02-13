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
	"errors"
	"fmt"
	"github.com/jmespath/go-jmespath"
	"io"
	"strings"
	"text/tabwriter"
)

type Format string

const (
	FormatText  = Format("text")
	FormatTable = Format("table")
	FormatJSON  = Format("json")
)

type Printer interface {
	Columns() []string
	Data() []map[string]string
	JSON() (io.Reader, error)
	Unwrap() interface{}
}

type Options struct {
	Format Format
	Query  string
}

func Print(p Printer, opts Options) (io.Reader, error) {
	switch opts.Format {
	case FormatText:
		return nil, errors.New("not implemented")
	case FormatJSON:
		return printJSON(p, opts)
	case FormatTable:
		return printTable(p, opts)
	default:
		return nil, errors.New("not implemented")
	}
}

func printJSON(p Printer, opts Options) (io.Reader, error) {
	if opts.Query == "" {
		return p.JSON()
	}

	data, err := json.Marshal(p.Unwrap())
	if err != nil {
		return nil, err
	}

	var x interface{}

	if err := json.Unmarshal(data, &x); err != nil {
		return nil, err
	}

	result, err := jmespath.Search(opts.Query, x)
	if err != nil {
		return nil, err
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(out), nil
}

func printTable(p Printer, _ Options) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, strings.Join(p.Columns(), "\t")); err != nil {
		return nil, err
	}

	for _, v := range p.Data() {
		row := make([]string, 0)

		for _, col := range p.Columns() {
			if v, ok := v[col]; ok {
				row = append(row, v)
			}
		}

		if _, err := fmt.Fprintln(tw, strings.Join(row, "\t")); err != nil {
			return nil, err
		}
	}

	if err := tw.Flush(); err != nil {
		return nil, err
	}

	return buf, nil
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}

	return "no"
}
