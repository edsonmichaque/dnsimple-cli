package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDNSSEC(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dnssec",
		Short: "Manage DNSEC status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		GroupID: "domains",
	}

	cmd.Flags().StringP("domain", "D", "", "Domain")
	cmd.Flags().BoolP("status", "s", false, "Get DNSSEC status")
	cmd.Flags().BoolP("enable", "y", false, "Enable DNSSEC status")
	cmd.Flags().BoolP("disable", "n", false, "Disable DNSSEC status")

	return cmd
}
