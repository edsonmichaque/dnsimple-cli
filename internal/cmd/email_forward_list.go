package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdEmailForwardList(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List email forwards",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
