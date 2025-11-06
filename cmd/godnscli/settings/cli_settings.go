package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/settings/dotenv"
)

// GetConfigDir returns the path to the CLI configuration directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".godns")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

func Init() {
	// Step 1: Set defaults FIRST
	// Keycloak/authentication settings for CLI
	viper.SetDefault(consts.KEYCLOAK_URL, "http://localhost:14101")
	viper.SetDefault(consts.KEYCLOAK_REALM, "godns")
	viper.SetDefault(consts.KEYCLOAK_CLI_CLIENT_ID, "godns-cli")

	// API settings
	viper.SetDefault(consts.HTTP_API_PORT, "14000")
	viper.SetDefault("api.url", "http://localhost:14000")

	// Development mode
	viper.SetDefault(consts.DEVELOPMENT, true)

	// Step 2: Load .env file (for backward compatibility)
	dotenv.LoadDotEnv()

	// Step 3: Read config file (this will override .env and defaults)
	configDir, err := GetConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not initialize config directory: %v\n", err)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(configDir)

		// Try to read existing config file (this will override defaults and .env)
		if err := viper.ReadInConfig(); err != nil {
			// It's okay if config doesn't exist yet, we'll create it on first save
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				fmt.Fprintf(os.Stderr, "Warning: Error reading config file: %v\n", err)
			}
		}
	}

	// Step 4: Allow environment variables to override everything
	viper.AutomaticEnv()
}

// SaveConfig writes the current configuration to the config file
func SaveConfig() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	return viper.WriteConfigAs(configPath)
}

// GetAPIURL returns the API URL from config
func GetAPIURL() string {
	return viper.GetString("api.url")
}

// GetAccessToken returns the cached access token
func GetAccessToken() (string, error) {
	// Try to get from pkg/auth token cache
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".godns")
	tokenPath := filepath.Join(configDir, "token.json")
	cleanPath := filepath.Clean(tokenPath)
	if !strings.HasPrefix(cleanPath, configDir+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid token path")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("not authenticated. Please run 'godnscli login' first")
		}
		return "", fmt.Errorf("failed to read token cache: %w", err)
	}

	var tokenCache struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(data, &tokenCache); err != nil {
		return "", fmt.Errorf("failed to parse token cache: %w", err)
	}

	if tokenCache.AccessToken == "" {
		return "", fmt.Errorf("no access token found. Please run 'godnscli login' first")
	}

	return tokenCache.AccessToken, nil
}
