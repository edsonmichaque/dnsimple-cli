package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdTransferAccept(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept",
		Short: "Accept transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
