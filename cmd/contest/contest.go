package contest

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "contest",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(checkCmd)
}
