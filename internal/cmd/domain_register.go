package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDomainRegister(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
