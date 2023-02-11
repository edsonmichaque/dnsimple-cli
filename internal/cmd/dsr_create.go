package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdDSRCreate(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create domain signed record",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
