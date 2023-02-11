package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDomain(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "domain",
		Short:   "Manage domains",
		GroupID: "domains",
	}

	cmd.AddCommand(NewCmdDomainCreate(opts))
	cmd.AddCommand(NewCmdDomainGet(opts))
	cmd.AddCommand(NewCmdDomainList(opts))
	cmd.AddCommand(NewCmdDomainDelete(opts))
	cmd.AddCommand(NewCmdDomainCheck(opts))
	cmd.AddCommand(NewCmdDomainRegister(opts))
	cmd.AddCommand(NewCmdDomainPrices(opts))
	cmd.AddCommand(NewCmdDomainRenew(opts))

	return cmd
}
