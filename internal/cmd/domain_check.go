package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDomainCheck(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
