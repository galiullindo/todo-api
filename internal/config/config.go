package config

import (
	"os"
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

func Load(path string) (*Config, error) {
	var configFile string
	switch env := os.Getenv("ENV"); env {
	case "test":
		configFile = "test.yaml"
	case "docker":
		configFile = "docker.yaml"
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
