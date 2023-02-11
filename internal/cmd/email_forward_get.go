package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdEmailForwardGet(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get email forward",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
