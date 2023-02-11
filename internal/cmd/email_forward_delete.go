package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdEmailForwardDelete(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete email forward",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
