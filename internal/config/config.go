package config

type Config struct {
	DSN  string
	Host string
	Port string
}

var config = &Config{
	DSN:  "/usr/local/share/production.db", // TODO: временно для docker
	Host: "0.0.0.0",                        // TODO: временно для docker
	Port: "8080",
}

func New() *Config {
	return config
}
