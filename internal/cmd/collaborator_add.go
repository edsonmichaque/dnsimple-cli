package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdCollaboratorAdd(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add collaborator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
