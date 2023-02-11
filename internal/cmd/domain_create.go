package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDomainCreate(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
