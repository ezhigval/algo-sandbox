package config

import "github.com/ezhigval/go-toolkit/config"

type Config struct {
	Port      string `env:"PORT" envDefault:"8084"`
	LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"json"`

	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`
}

func MustLoad() Config {
	return config.MustLoad[Config]()
}
