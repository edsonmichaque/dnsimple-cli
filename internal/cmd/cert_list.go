package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdCertList(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List certs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
