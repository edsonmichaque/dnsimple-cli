package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdCert(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cert",
		Short:   "Manage certs",
		GroupID: "certs",
	}

	cmd.AddCommand(NewCmdCertGet(opts))
	cmd.AddCommand(NewCmdCertList(opts))
	cmd.AddCommand(NewCmdCertDownload(opts))

	return cmd
}
