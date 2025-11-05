package cmd

import (
	"fmt"
	"time"

	"github.com/rogerwesterbo/godns/pkg/auth"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long: `Display the current authentication status including token validity.

This command checks if you are currently logged in and shows:
- Whether a valid access token exists
- When the token expires
- Whether the token needs to be refreshed

Example:
  godnscli status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cache, err := auth.LoadTokenCache()
		if err != nil {
			return fmt.Errorf("failed to load token cache: %w", err)
		}

		if cache == nil {
			fmt.Println("❌ Not logged in")
			fmt.Println("\nRun 'godnscli login' to authenticate")
			return nil
		}

		fmt.Println("🔐 Authentication Status")
		fmt.Println("========================")

		if cache.IsTokenValid() {
			fmt.Println("Status: ✅ Logged in")
			fmt.Printf("Token Type: %s\n", cache.TokenType)

			timeUntilExpiry := time.Until(cache.ExpiresAt)
			fmt.Printf("Expires: %s (in %s)\n",
				cache.ExpiresAt.Format("2006-01-02 15:04:05"),
				formatDuration(timeUntilExpiry))

			if cache.RefreshToken != "" {
				fmt.Println("Refresh Token: Available")
			}
		} else {
			fmt.Println("Status: ⚠️  Token expired")
			fmt.Printf("Expired: %s\n", cache.ExpiresAt.Format("2006-01-02 15:04:05"))

			if cache.RefreshToken != "" {
				fmt.Println("\nToken will be automatically refreshed on next API call")
			} else {
				fmt.Println("\nRun 'godnscli login' to re-authenticate")
			}
		}

		return nil
	},
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		return "expired"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 24 {
		days := hours / 24
		hours = hours % 24
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}

	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	seconds := int(d.Seconds())
	return fmt.Sprintf("%ds", seconds)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
