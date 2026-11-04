package config

type Config struct {
	DBPath string
	Host   string
	Port   string
}

var config = &Config{
	DBPath: "todo.db",
	Host:   "localhost",
	Port:   "8080",
}

func New() *Config {
	return config
}
