package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdLetsEncryptIssue(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue letsencrypt certificate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
