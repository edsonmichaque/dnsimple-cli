package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdPushReject(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject",
		Short: "Reject push",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
