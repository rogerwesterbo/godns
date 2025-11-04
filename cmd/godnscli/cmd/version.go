package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version information",
	Long:    `Print the version, commit, and build date of godnscli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("godnscli version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built:  %s\n", date)
	},
}

func initVersionCommand() {
	rootCmd.AddCommand(versionCmd)
}
