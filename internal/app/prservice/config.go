package prservice

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	Service  ServiceConfig
	Database DatabaseConfig
}

type ServiceConfig struct {
	ServerConfig
	LogLevel string `env:"LOG_LEVEL,required"`
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT,required"`
}

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST,required"`
	Port     string `env:"DATABASE_PORT,required"`
	User     string `env:"DATABASE_USER,required"`
	Password string `env:"DATABASE_PASSWORD,required"`
	Name     string `env:"DATABASE_NAME,required"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
