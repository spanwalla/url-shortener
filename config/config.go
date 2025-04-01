package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App        `yaml:"app"`
		HTTP       `yaml:"http"`
		Log        `yaml:"log"`
		PG         `yaml:"postgres"`
		Alias      `yaml:"alias"`
		Shortening `yaml:"shortening"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" env:"PG_URL"`
	}

	Alias struct {
		Alphabet string `env-required:"true" yaml:"alphabet" env:"ALIAS_ALPHABET"`
		Length   int    `env-required:"true" yaml:"length" env:"ALIAS_LENGTH"`
	}

	Shortening struct {
		Attempts int `env-required:"true" yaml:"attempts" env:"SHORTENING_ATTEMPTS"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.ReadConfig: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.UpdateEnv: %w", err)
	}

	return cfg, nil
}
