package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdPushInit(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init push",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
