package cmd

import (
	"io"

	"context"
	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type CmdOpt struct {
	Stdout io.Writer
	Stdin  io.Reader
	Stderr io.Writer

	Client func(string, string) *dnsimple.Client
}

func NewCmdOpt() *CmdOpt {
	buildClient := func(url, token string) *dnsimple.Client {
		client := dnsimple.NewClient(dnsimple.StaticTokenHTTPClient(context.Background(), token))

		if url != "" {
			client.BaseURL = url
		}

		return client
	}

	return &CmdOpt{
		Client: buildClient,
	}
}
