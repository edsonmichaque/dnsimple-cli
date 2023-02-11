package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDomainRenew(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew",
		Short: "Renew domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
