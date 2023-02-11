package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdTLDExtendedAttr(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ext",
		Short: "Get TLD extended attributes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
