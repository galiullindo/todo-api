package config

import (
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	URI string `yaml:"uri"`
}

func Load(path string, env Env) (*Config, error) {
	var configFile string
	switch env {
	case EnvProd:
		configFile = "config.prod.yaml"
	case EnvTest:
		configFile = "config.test.yaml"
	case EnvDocker:
		configFile = "config.docker.yaml"
	default:
		configFile = "config.yaml"
	}
	configPath := filepath.Clean(filepath.Join(path, configFile))

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
