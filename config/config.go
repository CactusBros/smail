package config

import "github.com/caarlos0/env/v11"

type Config struct {
	HTTP HTTPConfig
	SMTP SMTPConfig
}

type HTTPConfig struct {
	Port int `env:"HTTP_PORT"`
}

type SMTPConfig struct {
	Host     string `env:"SMTP_HOST"`
	Port     int    `env:"SMTP_PORT"`
	From     string `env:"SMTP_FROM"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
}

func MustReadConfigFromEnv() Config {
	cfg, err := ReadConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return cfg
}

func ReadConfigFromEnv() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
