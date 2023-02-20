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
	"io"
	"strings"
	"text/tabwriter"

	"github.com/jmespath/go-jmespath"
	"gopkg.in/yaml.v3"
)

type Format string

const (
	FormatText  = Format("text")
	FormatTable = Format("table")
	FormatJSON  = Format("json")
	FormatYAML  = Format("yaml")
)

type Options struct {
	Format Format
	Query  string
}

type YAMLPrinter interface {
	PrintYAML(opts *Options) (io.Reader, error)
}

type JSONPrinter interface {
	PrintJSON(opts *Options) (io.Reader, error)
}

type TablePrinter interface {
	PrintTable(opts *Options) (io.Reader, error)
}

type TextPrinter interface {
	PrintText(opts *Options) (io.Reader, error)
}

func Print(data interface{}, opts *Options) (io.Reader, error) {
	if opts.Format == FormatJSON {
		if printer, ok := data.(JSONPrinter); ok {
			return printer.PrintJSON(opts)
		}

		return nil, errors.New("json printer is not implemented")
	}

	if opts.Format == FormatYAML {
		if printer, ok := data.(YAMLPrinter); ok {
			return printer.PrintYAML(opts)
		}

		return nil, errors.New("yaml printer is not implemented")
	}

	if opts.Format == FormatTable {
		if printer, ok := data.(TablePrinter); ok {
			return printer.PrintTable(opts)
		}

		return nil, errors.New("table printer is not implemented")
	}

	if opts.Format == FormatText {
		if printer, ok := data.(TextPrinter); ok {
			return printer.PrintText(opts)
		}

		return nil, errors.New("table printer is not implemented")
	}

	return nil, errors.New("invalid printer")
}

func printJSON(j jsonPrinter, opts *Options) (io.Reader, error) {
	data, err := j.toJSON(opts)
	if err != nil {
		return nil, err
	}

	var result interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if opts.Query != "" {
		result, err = jmespath.Search(opts.Query, result)
		if err != nil {
			return nil, err
		}
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(out), nil
}

func printYAML(j jsonPrinter, opts *Options) (io.Reader, error) {
	data, err := j.toJSON(opts)
	if err != nil {
		return nil, err
	}

	var result interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if opts.Query != "" {
		result, err = jmespath.Search(opts.Query, result)
		if err != nil {
			return nil, err
		}
	}

	out, err := yaml.Marshal(result)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(out), nil
}

type jsonPrinter interface {
	toJSON(opts *Options) ([]byte, error)
}

type tablePrinter interface {
	printHeader() []string
	printRows() []map[string]string
}

func printTable(t tablePrinter) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, strings.Join(t.printHeader(), "\t")); err != nil {
		return nil, err
	}

	for _, v := range t.printRows() {
		row := make([]string, 0)

		for _, col := range t.printHeader() {
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
