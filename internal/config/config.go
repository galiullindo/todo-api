package config

type Config struct {
	DatabaseURI string
	Host        string
	Port        string
}

var config = &Config{
	DatabaseURI: "production.db",
	Host:        "localhost",
	Port:        "8080",
}

func New() *Config {
	return config
}
