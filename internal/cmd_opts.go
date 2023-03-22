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

package internal

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/spf13/cobra"
)

func NewCommandOpts() *CmdOpts {
	return &CmdOpts{
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
		Stdout: os.Stdout,
		Client: buildClient,
	}
}

type CmdOpts struct {
	Stdout io.Writer
	Stdin  io.Reader
	Stderr io.Writer

	Client func(string, string) *dnsimple.Client
}

func (c CmdOpts) Validate() error {
	if c.Client == nil {
		return errors.New("invalid client builder")
	}

	return nil
}

func (c CmdOpts) BuildClient(url, token string) *dnsimple.Client {
	client := c.Client

	if client == nil {
		client = buildClient
	}

	return client(url, token)
}

func buildClient(url, token string) *dnsimple.Client {
	client := dnsimple.NewClient(dnsimple.StaticTokenHTTPClient(context.Background(), token))

	if url != "" {
		client.BaseURL = url
	}

	return client
}

func SetupIO(cmd *cobra.Command, opts *CmdOpts) {
	if opts.Stdout != nil {
		cmd.SetOut(opts.Stdout)
	}

	if opts.Stdin != nil {
		cmd.SetIn(opts.Stdin)
	}

	if opts.Stderr != nil {
		cmd.SetErr(opts.Stderr)
	}
}
