package config

import "errors"

var ErrUnknownEnviroment = errors.New("unknown enviroment")

type Env string

const (
	EnvDefault Env = "default"
	EnvProd    Env = "prod"
	EnvTest    Env = "test"
	EnvDocker  Env = "docker"
)

func ParseEnv(s string) (Env, error) {
	switch s {
	case "default":
		return EnvDefault, nil
	case "prod":
		return EnvProd, nil
	case "test":
		return EnvTest, nil
	case "docker":
		return EnvDocker, nil
	default:
		return "", ErrUnknownEnviroment
	}
}
