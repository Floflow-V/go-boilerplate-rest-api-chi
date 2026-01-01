package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Api      ApiConfig      `envPrefix:"API_"`
	Log      LogConfig      `envPrefix:"LOG_"`
	Database DatabaseConfig `envPrefix:"DATABASE_"`
}

type ApiConfig struct {
	Environement string `env:"ENVIRONEMENT,required,notEmpty"`
	Host         string `env:"HOST,required,notEmpty"`
	Port         int    `env:"PORT,required,notEmpty"`
}

type LogConfig struct {
	Level  string `env:"LEVEL,required,notEmpty"`
	Format string `env:"FORMAT,required,notEmpty"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST,required,notEmpty"`
	Port     int    `env:"PORT,required,notEmpty"`
	User     string `env:"USER,required,notEmpty"`
	Password string `env:"PASSWORD,required,notEmpty"`
	Name     string `env:"NAME,required,notEmpty"`
	LogLevel string `env:"LOG_LEVEL,required,notEmpty"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
