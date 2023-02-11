package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdTransfer(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transfer",
		Short:   "Manage transfers",
		GroupID: "domains",
	}

	cmd.AddCommand(NewCmdTransferCreate(opts))
	cmd.AddCommand(NewCmdTransferGet(opts))
	cmd.AddCommand(NewCmdTransferCancel(opts))
	cmd.AddCommand(NewCmdTransferAccept(opts))

	return cmd
}
