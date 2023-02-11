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

package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	modName := "github.com/edsonmichaque/dnsimple-cli"

	ldflags := []string{
		fmt.Sprintf("-X %s/internal/build.Version=%s", modName, version()),
		fmt.Sprintf("-X %s/internal/build.Date=%s", modName, date()),
	}

	command := []string{
		"go",
		"build",
		strings.Join(ldflags, " "),
		"-o",
		"bin/dnsimple",
		"cmd/dnsimple/main.go",
	}

	fmt.Println(command)
}

func version() string {
	// check DNSIMPLE_CLI_VERSION
	// run git describe --tags
	// run git rev-parse --shord HEAD

	if ver := os.Getenv("DNSIMPLE_CLI_VERSION"); ver != "" {
		return ver
	}

	return "dev"
}

func date() string {
	return time.Now().Format("2006-01-02")
}
