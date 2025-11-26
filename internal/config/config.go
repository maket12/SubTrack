package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPAddress string `env:"HTTP_ADDRESS" envDefault:":8080"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`

	DatabaseDSN string `env:"DATABASE_DSN,required"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
