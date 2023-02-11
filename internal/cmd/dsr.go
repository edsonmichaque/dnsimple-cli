package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDSR(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dsr",
		Short:   "Manage domain signed records",
		GroupID: "domains",
	}

	cmd.AddCommand(NewCmdDSRCreate(opts))
	cmd.AddCommand(NewCmdDSRGet(opts))
	cmd.AddCommand(NewCmdDSRList(opts))
	cmd.AddCommand(NewCmdDSRDelete(opts))

	return cmd
}
