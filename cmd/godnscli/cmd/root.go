package cmd

import (
	"fmt"
	"os"

	"github.com/rogerwesterbo/godns/cmd/godnscli/settings"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godnscli",
	Short: "GoDNS CLI - Test and manage your GoDNS server",
	Long: `GoDNS CLI is a command-line tool to test and interact with your GoDNS server.
	
It provides commands to:
- Query DNS records (A, AAAA, MX, NS, etc.)
- Test DNS server functionality
- Manage DNS records in Valkey
- Check server health`,
}

func Execute() {
	// Initialize settings before executing commands
	settings.Init()

	// Initialize all command flags
	initRootFlags()
	initQueryCommand()
	initHealthCommand()
	initDiscoverCommand()
	initTestCommand()
	initVersionCommand()
	initExportCommand()
	// Config command initializes itself in init()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initRootFlags() {
	// Global flags can be added here
	rootCmd.PersistentFlags().StringP("server", "s", "localhost:53", "DNS server address")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
