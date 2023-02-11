package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdPush(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "push",
		Short:   "Manage pushes",
		GroupID: "domains",
	}

	cmd.AddCommand(NewCmdPushInit(opts))
	cmd.AddCommand(NewCmdPushAccept(opts))
	cmd.AddCommand(NewCmdPushReject(opts))
	cmd.AddCommand(NewCmdPushList(opts))

	return cmd
}
