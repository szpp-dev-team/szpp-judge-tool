package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/szpp-dev-team/szpp-judge-tool/cmd/contest"
	"github.com/szpp-dev-team/szpp-judge-tool/cmd/task"
)

var rootCmd = &cobra.Command{
	Use: "szpp-judge-tool",
}

func init() {
	rootCmd.AddCommand(task.Cmd, contest.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
