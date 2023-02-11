package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdPushAccept(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept",
		Short: "Accept push",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
