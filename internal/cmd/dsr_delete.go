package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDSRDelete(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete domain signed record",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
