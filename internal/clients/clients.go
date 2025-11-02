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

	v1Valkeyclient, err := v1valkeyclient.NewV1ValkeyClient(&valkeyOpts)
	if err != nil {
		vlog.Fatalf("Failed to initialize Valkey client: %v", err)
	}
	V1ValkeyClient = v1Valkeyclient
}
