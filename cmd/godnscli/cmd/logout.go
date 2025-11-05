package cmd

import (
	"github.com/rogerwesterbo/godns/pkg/auth"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from GoDNS API",
	Long: `Logout from GoDNS by removing cached authentication tokens.

After logout, you will need to run 'godnscli login' again to authenticate.

Example:
  godnscli logout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return auth.Logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
