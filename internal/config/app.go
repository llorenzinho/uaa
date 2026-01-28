package config

import (
	"strings"
	"sync"

	"github.com/llorenzinho/goauth/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConfig struct {
	ServerConfig `mapstructure:"server" validate:"required"`
	DBConfig     `mapstructure:"database" validate:"required"`
}

var once sync.Once
var cfg AppConfig

func NewAppConfig() *AppConfig {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")

		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // this is useful e.g. want to use . in Get() calls, but environmental variables to use _ delimiters (e.g. app.port -> APP_PORT)

		if err := viper.ReadInConfig(); err != nil {
			log.Get().Fatal("failed to read config file", zap.Error(err))
		}

		if err := viper.Unmarshal(&cfg); err != nil {
			log.Get().Fatal("failed to unmarshal config", zap.Error(err))
		}
	})

	return &cfg
}
