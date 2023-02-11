package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDomainDelete(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
