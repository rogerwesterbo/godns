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

	// GoDNS Application ports
	viper.SetDefault(consts.DNS_SERVER_PORT, ":53")
	viper.SetDefault(consts.DNS_SERVER_LIVENESS_PROBE_PORT, ":14003")
	viper.SetDefault(consts.DNS_SERVER_READYNESS_PROBE_PORT, ":14004")
	viper.SetDefault(consts.DNS_ENABLE_HTTP_API, true)
	viper.SetDefault(consts.DNS_ENABLE_ALLOWED_LANS_CHECK, false) // Default off for development
	viper.SetDefault(consts.HTTP_API_PORT, ":8080")
	viper.SetDefault(consts.HTTP_API_READINESS_PROBE_PORT, ":8081")
	viper.SetDefault(consts.HTTP_API_LIVENESS_PROBE_PORT, ":8082")

	// DNS Cache settings
	viper.SetDefault(consts.DNS_CACHE_ENABLED, true)
	viper.SetDefault(consts.DNS_CACHE_SIZE, 10000)
	viper.SetDefault(consts.DNS_CACHE_TTL_SECONDS, 300) // 5 minutes

	// DNS Rate Limiting settings
	viper.SetDefault(consts.DNS_RATE_LIMIT_ENABLED, true)
	viper.SetDefault(consts.DNS_RATE_LIMIT_QUERIES_PER_SEC, 100) // 100 queries per second per IP
	viper.SetDefault(consts.DNS_RATE_LIMIT_BURST, 200)           // Allow bursts up to 200

	// DNS Load Balancing settings
	viper.SetDefault(consts.DNS_LOAD_BALANCER_ENABLED, true)
	viper.SetDefault(consts.DNS_LOAD_BALANCER_STRATEGY, "round-robin")

	// DNS Health Check settings
	viper.SetDefault(consts.DNS_HEALTH_CHECK_ENABLED, false) // Default off, enable when needed
	viper.SetDefault(consts.DNS_HEALTH_CHECK_INTERVAL_SEC, 30)
	viper.SetDefault(consts.DNS_HEALTH_CHECK_TIMEOUT_SEC, 5)

	// DNS Query Logging settings
	viper.SetDefault(consts.DNS_QUERY_LOG_ENABLED, true)
	viper.SetDefault(consts.DNS_QUERY_LOG_BUFFER_SIZE, 1000)
	viper.SetDefault(consts.DNS_QUERY_LOG_FLUSH_INTERVAL, "1m")

	// Metrics settings
	viper.SetDefault(consts.METRICS_ENABLED, true)
	viper.SetDefault(consts.METRICS_PORT, ":9090")

	viper.SetDefault(consts.VALKEY_HOST, "localhost")
	viper.SetDefault(consts.VALKEY_PORT, "6379")
	viper.SetDefault(consts.VALKEY_TOKEN, "")
	viper.SetDefault(consts.VALKEY_MAX_RETRIES, 3)
	viper.SetDefault(consts.VALKEY_INITIAL_RETRY_DELAY_MS, 100)

	// Authentication settings
	viper.SetDefault(consts.AUTH_ENABLED, true)
	viper.SetDefault(consts.KEYCLOAK_URL, "http://localhost:14101")
	viper.SetDefault(consts.KEYCLOAK_REALM, "godns")
	viper.SetDefault(consts.KEYCLOAK_API_CLIENT_ID, "godns-api")
	viper.SetDefault(consts.KEYCLOAK_CLI_CLIENT_ID, "godns-cli")

	viper.AutomaticEnv()

	dotenv.LoadDotEnv()
}
