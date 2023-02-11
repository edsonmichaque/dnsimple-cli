package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdTransferCancel(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
