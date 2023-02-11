package cmd

import "io"

type CmdOpt struct {
	Stdout io.Writer
	Stdin  io.Reader
	Stderr io.Writer
}
