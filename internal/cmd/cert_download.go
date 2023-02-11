package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdCertDownload(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download cert",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
