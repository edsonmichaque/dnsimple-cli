package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdWhoami(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Check identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		GroupID: "identity",
	}

	return cmd
}
