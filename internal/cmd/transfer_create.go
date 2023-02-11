package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdTransferCreate(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
