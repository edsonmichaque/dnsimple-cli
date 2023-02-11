package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDomainPrices(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prices",
		Short: "Get domain prices",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
