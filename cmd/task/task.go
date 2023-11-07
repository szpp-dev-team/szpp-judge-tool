package task

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/szpp-dev-team/szpp-judge-tool/internal/task"
)

var Cmd = &cobra.Command{
	Use: "task",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var checkCmd = &cobra.Command{
	Use: "check",
	RunE: func(_ *cobra.Command, args []string) error {
		dirPath, err := os.Getwd()
		if err != nil {
			return err
		}
		if len(args) > 0 {
			dirPath = args[0]
		}

		controller, err := task.Load(dirPath)
		if err != nil {
			return err
		}
		defer controller.Cleanup()
		if err := controller.Validate(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(checkCmd)
}
