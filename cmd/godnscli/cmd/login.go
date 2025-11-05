package cmd

import (
	"context"

	"github.com/rogerwesterbo/godns/pkg/auth"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with GoDNS API using Keycloak",
	Long: `Login to GoDNS using OAuth2 device flow authentication.

This command will:
1. Request a device code from Keycloak
2. Display a URL and code for you to authorize
3. Wait for you to complete authentication in your browser
4. Save the access token locally for future use

Example:
  godnscli login`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return auth.DeviceCodeLogin(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
