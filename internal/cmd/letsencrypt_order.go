package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdLetsEncryptOrder(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order",
		Short: "Order letsencrypt certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
