package clients

import (
	"time"

	"github.com/rogerwesterbo/godns/pkg/clients/v1valkeyclient"
	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/rogerwesterbo/godns/pkg/options/valkeyoptions"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

var (
	V1ValkeyClient *v1valkeyclient.V1ValkeyClient
)

func Init() {
	valkeyOpts := valkeyoptions.ValkeyOptions{
		Host:              viper.GetString(consts.VALKEY_HOST),
		Port:              viper.GetString(consts.VALKEY_PORT),
		Username:          viper.GetString(consts.VALKEY_USERNAME),
		APIToken:          viper.GetString(consts.VALKEY_TOKEN),
		TimeoutSec:        30,
		MaxRetries:        viper.GetInt(consts.VALKEY_MAX_RETRIES),
		InitialRetryDelay: time.Duration(viper.GetInt(consts.VALKEY_INITIAL_RETRY_DELAY_MS)) * time.Millisecond,
	}

	// Log configuration (without sensitive data) for troubleshooting
	vlog.Infof("Initializing Valkey client with host=%s, port=%s, username=%s",
		valkeyOpts.Host, valkeyOpts.Port, valkeyOpts.Username)

	// Check for common configuration mistakes
	if valkeyOpts.Username != "" && valkeyOpts.APIToken == "" {
		vlog.Warnf("VALKEY_USERNAME is set but VALKEY_TOKEN is empty. This may cause authentication failures.")
		vlog.Warnf("Please ensure VALKEY_TOKEN environment variable is set to match your Valkey ACL configuration.")
	}

	v1Valkeyclient, err := v1valkeyclient.NewV1ValkeyClient(&valkeyOpts)
	if err != nil {
		vlog.Errorf("Failed to initialize Valkey client: %v", err)
		vlog.Errorf("Troubleshooting tips:")
		vlog.Errorf("  1. Verify VALKEY_HOST=%s and VALKEY_PORT=%s are correct", valkeyOpts.Host, valkeyOpts.Port)
		vlog.Errorf("  2. Verify VALKEY_USERNAME=%s matches a user in hack/valkey/users.acl", valkeyOpts.Username)
		vlog.Errorf("  3. Verify VALKEY_TOKEN is set and matches the password in hack/valkey/users.acl")
		vlog.Errorf("  4. If Valkey data was persisted with old credentials, clear it: rm -rf hack/data/valkey/*")
		vlog.Fatalf("Cannot start without Valkey connection")
	}
	V1ValkeyClient = v1Valkeyclient
}
