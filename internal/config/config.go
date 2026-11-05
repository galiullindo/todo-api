package config

type Config struct {
	DSN  string
	Host string
	Port string
}

var config = &Config{
	DSN:  "../data/todo.db",
	Host: "localhost",
	Port: "8080",
}

func New() *Config {
	return config
}
