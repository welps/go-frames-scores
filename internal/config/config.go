package config

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/welps/go-frames-scores/internal/constants"
)

type Config struct {
	Environment        constants.Environment `mapstructure:"ENVIRONMENT"`
	Port               int                   `mapstructure:"PORT"`
	GracefulShutdownMS int                   `mapstructure:"GRACEFUL_SHUTDOWN_MS"`
	PublicURL          string                `mapstructure:"PUBLIC_URL"`

	HTTPClientSettings HTTPClientSettings `mapstructure:",squash"`
	SportsAPIConfig    SportsAPIConfig    `mapstructure:",squash"`
}

type HTTPClientSettings struct {
	MaxIdleConns        int `mapstructure:"MAX_IDLE_CONNS"`
	MaxIdleConnsPerHost int `mapstructure:"MAX_IDLE_CONNS_PER_HOST"`
	RequestTimeoutMS    int `mapstructure:"REQUEST_TIMEOUT_MS"`
}

type SportsAPIConfig struct {
	Host   string `mapstructure:"SPORTS_API_HOST"`
	APIKey string `mapstructure:"SPORTS_API_KEY"`
}

func InitConfig() Config {
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("GRACEFUL_SHUTDOWN_MS", (10 * time.Second).Milliseconds())
	viper.SetDefault("PUBLIC_URL", "http://localhost:8080")

	viper.SetDefault("MAX_IDLE_CONNS", 100)
	viper.SetDefault("MAX_IDLE_CONNS_PER_HOST", 50)
	viper.SetDefault("REQUEST_TIMEOUT_MS", (30 * time.Second).Milliseconds())

	viper.SetDefault("SPORTS_API_HOST", "https://sportscore1.p.rapidapi.com")
	viper.SetDefault("SPORTS_API_KEY", "")

	viper.AutomaticEnv()

	config := Config{}
	_ = viper.Unmarshal(
		&config,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
				mapstructure.StringToSliceHookFunc(","),
			),
		),
	)

	return config
}
