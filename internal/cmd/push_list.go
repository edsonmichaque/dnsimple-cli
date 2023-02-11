package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdPushList(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
