package task

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "task",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(checkCmd, uploadCmd)
}
