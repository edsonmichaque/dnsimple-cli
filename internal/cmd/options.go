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
	"context"
	"errors"
	"io"
	"os"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

func NewOptions() (*Options, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Options{
		Stdin:         os.Stdin,
		Stderr:        os.Stderr,
		Stdout:        os.Stdout,
		WorkDir:       wd,
		ClientBuilder: buildClient,
	}, nil
}

type Options struct {
	Stdout        io.Writer
	Stdin         io.Reader
	Stderr        io.Writer
	WorkDir       string
	ClientBuilder func(string, string) *dnsimple.Client
}

func (c Options) Validate() error {
	if c.ClientBuilder == nil {
		return errors.New("invalid client builder")
	}

	return nil
}

func (c Options) BuildClient(url, token string) *dnsimple.Client {
	clientBuilder := c.ClientBuilder

	if clientBuilder == nil {
		clientBuilder = buildClient
	}

	return clientBuilder(url, token)
}

func buildClient(url, token string) *dnsimple.Client {
	client := dnsimple.NewClient(dnsimple.StaticTokenHTTPClient(context.Background(), token))

	if url != "" {
		client.BaseURL = url
	}

	return client
}
