package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdTLD(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tld",
		Short:   "Manage tlds",
		GroupID: "domains",
	}

	cmd.AddCommand(NewCmdTLDGet(opts))
	cmd.AddCommand(NewCmdTLDList(opts))
	cmd.AddCommand(NewCmdTLDExtendedAttr(opts))

	return cmd
}
