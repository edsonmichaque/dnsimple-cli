package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdEmailForwardCreate(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create email forward",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
