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
