package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rogerwesterbo/godns/cmd/godnscli/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long: `View and manage GoDNS CLI configuration.

The configuration is stored at ~/.godns/config.yaml and can include:
- API endpoint URLs
- Keycloak authentication settings
- Development mode flags

Examples:
  godnscli config show          # Show current configuration
  godnscli config path          # Show config file path
  godnscli config set api.url http://localhost:14000
  godnscli config get api.url`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Current Configuration:")
		fmt.Println("======================")

		allSettings := viper.AllSettings()
		for key, value := range allSettings {
			fmt.Printf("%s: %v\n", key, value)
		}

		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := settings.GetConfigDir()
		if err != nil {
			return err
		}

		configPath := filepath.Join(configDir, "config.yaml")
		fmt.Println(configPath)

		// Check if file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println("(file does not exist yet - will be created on first use)")
		}

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		viper.Set(key, value)

		if err := settings.SaveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ Set %s = %s\n", key, value)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := viper.Get(key)

		if value == nil {
			fmt.Printf("%s: (not set)\n", key)
		} else {
			fmt.Printf("%s: %v\n", key, value)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
}
