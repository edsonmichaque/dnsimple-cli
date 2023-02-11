package cmd

import (
	"github.com/spf13/cobra"
)

func NewCmdCollaborator(opts *CmdOpt) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collaborator",
		Short:   "Manage collaborators",
		GroupID: "domains",
	}

	cmd.AddCommand(NewCmdCollaboratorAdd(opts))
	cmd.AddCommand(NewCmdCollaboratorList(opts))

	return cmd
}
