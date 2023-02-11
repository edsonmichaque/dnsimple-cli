package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdLetsEncryptIssueRenewal(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue-renewal",
		Short: "Issue letsencrypt renewal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
