package settings

import (
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/settings/dotenv"
)

func Init() {
	viper.AutomaticEnv()

	dotenv.LoadDotEnv()

}
