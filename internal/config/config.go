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

package config

import (
	"errors"

	"github.com/spf13/viper"
)

const (
	BaseURLProduction = "https://api.dnsimple.com"
	BaseURLSandbox    = "https://api.sandbox.dnsimple.com"
)

func New() (*Config, error) {
	return NewWithValidation(true)
}

func NewWithValidation(validate bool) (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if validate {
		if err := cfg.Validate(); err != nil {
			return nil, err
		}
	}

	if cfg.Sandbox {
		cfg.BaseURL = BaseURLSandbox
	}

	return &cfg, nil
}

type Config struct {
	Account     string `mapstructure:"account"`
	Sandbox     bool   `mapstructure:"sandbox"`
	AccessToken string `mapstructure:"access-token"`
	BaseURL     string `mapstructure:"base-url"`
}

func (c Config) Validate() error {
	if c.Account == "" {
		return errors.New("account id is required")
	}

	if c.AccessToken == "" {
		return errors.New("access token is required")
	}

	return nil
}
