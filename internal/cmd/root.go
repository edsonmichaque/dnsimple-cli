package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdRoot(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use: "dnsimple",
	}

	cmd.AddGroup(&cobra.Group{
		ID:    "domains",
		Title: "Domains commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    "identity",
		Title: "Identity commands",
	})

	cmd.AddGroup(&cobra.Group{
		ID:    "certs",
		Title: "Certificates commands",
	})

	cmd.AddCommand(NewCmdDomain(opts))
	cmd.AddCommand(NewCmdDSR(opts))
	cmd.AddCommand(NewCmdDNSSEC(opts))
	cmd.AddCommand(NewCmdCollaborator(opts))
	cmd.AddCommand(NewCmdEmailForward(opts))
	cmd.AddCommand(NewCmdWhoami(opts))
	cmd.AddCommand(NewCmdPush(opts))
	cmd.AddCommand(NewCmdCert(opts))
	cmd.AddCommand(NewCmdTLD(opts))
	cmd.AddCommand(NewCmdLetsEncrypt(opts))
	cmd.AddCommand(NewCmdTransfer(opts))

	return cmd
}
