package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDSRList(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List domains domain signed records",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
