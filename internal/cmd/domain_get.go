package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDomainGet(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
