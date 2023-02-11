package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdEmailForward(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "email-forward",
		Short: "Manage email forwards",
		GroupID: "domains",
	}

	cmd.AddCommand(NewCmdEmailForwardCreate(opts))
	cmd.AddCommand(NewCmdEmailForwardGet(opts))
	cmd.AddCommand(NewCmdEmailForwardList(opts))
	cmd.AddCommand(NewCmdEmailForwardDelete(opts))

	return cmd
}
