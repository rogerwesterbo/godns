package settings

import (
	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/settings/dotenv"
)

func Init() {
	viper.SetDefault(consts.DEVELOPMENT, false)
	viper.SetDefault(consts.LOG_LEVEL, "info")
	viper.SetDefault(consts.LOG_JSON, true)
	viper.SetDefault(consts.LOG_ADD_CALLER, true)
	viper.SetDefault(consts.LOG_DISABLE_STACKTRACE, false)
	viper.SetDefault(consts.LOG_COLORIZE_LINE, false)
	viper.SetDefault(consts.LOG_UNESCAPE_MULTILINE, false)

	viper.SetDefault(consts.DNS_SERVER_PORT, ":53")
	viper.SetDefault(consts.DNS_SERVER_LIVENESS_PROBE_PORT, ":8080")
	viper.SetDefault(consts.DNS_SERVER_READYNESS_PROBE_PORT, ":8081")

	viper.SetDefault(consts.VALKEY_HOST, "localhost")
	viper.SetDefault(consts.VALKEY_PORT, "6379")
	viper.SetDefault(consts.VALKEY_TOKEN, "")
	viper.SetDefault(consts.VALKEY_MAX_RETRIES, 3)
	viper.SetDefault(consts.VALKEY_INITIAL_RETRY_DELAY_MS, 100)

	viper.AutomaticEnv()

	dotenv.LoadDotEnv()
}
