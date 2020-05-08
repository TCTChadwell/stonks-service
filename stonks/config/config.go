package config

import (
	"github.com/kelseyhightower/envconfig"
)

var Config StonkServiceConfig

type StonkServiceConfig struct {
	AppPort       string `envconfig:"APP_PORT" default:"9001"`
	RedisHost     string `envconfig:"STONK_REDIS_HOST" required:"true"`
	RedisPassword string `envconfig:"STONK_REDIS_PASS"`
	// to-add later
	NewsApiHost   string `envconfig:"NEWS_API_HOST"`
	TradingApiKey string `envconfig:"TRADING_API_KEY" required:"true"`
}

func CreateConfig() error {
	var s StonkServiceConfig
	err := envconfig.Process("", &s)
	if err != nil {
		return err
	}
	Config = s
	return nil
}
