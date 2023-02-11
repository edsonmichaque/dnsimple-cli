package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdLetsEncryptOrderRenewal(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order-renewal",
		Short: "Order letsencrypt renewal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
