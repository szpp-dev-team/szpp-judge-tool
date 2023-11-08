package contest

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/szpp-dev-team/szpp-judge-tool/internal/contest"
)

var checkCmd = &cobra.Command{
	Use: "check",
	RunE: func(cmd *cobra.Command, args []string) error {
		dirPath, err := os.Getwd()
		if err != nil {
			return err
		}

		controller, err := contest.Load(dirPath)
		if err != nil {
			return err
		}
		return controller.Validate()
	},
}
