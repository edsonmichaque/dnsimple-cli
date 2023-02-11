package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdLetsEncrypt(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "letsencrypt",
		Short:   "Manage letsencrypt certificates",
		GroupID: "certs",
	}

	cmd.AddCommand(NewCmdLetsEncryptOrder(opts))
	cmd.AddCommand(NewCmdLetsEncryptIssue(opts))
	cmd.AddCommand(NewCmdLetsEncryptOrderRenewal(opts))
	cmd.AddCommand(NewCmdLetsEncryptIssueRenewal(opts))

	return cmd
}
