package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Scheme             string `env:"SCHEME" envDefault:"http"`
	Host               string `env:"HOST" envDefault:"127.0.0.1"`
	Port               int    `env:"PORT" envDefault:"9091"`
	Path               string `env:"PATH" envDefault:"/transmission/rpc"`
	User               string `env:"USER"`
	Password           string `env:"PASSWORD"`
	OtelEndpoint       string `env:"OTEL_ENDPOINT" envDefault:"http://localhost:4318"`
	FetchInterval      int    `env:"FETCH_INTERVAL" envDefault:"10"`
	Debug              bool   `env:"DEBUG" envDefault:"false"`
	SpaceCheckPath     string `env:"SPACE_CHECK_PATH" envDefault:""`
	ErrorCheckInterval int    `env:"ERROR_CHECK_INTERVAL" envDefault:"60"`
}

var (
	log = logrus.WithFields(logrus.Fields{
		"prefix": "config",
	})
	C Config
)

func Load() (Config, error) {
	var cfg Config
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: "TRANSMISSION_"}); err != nil {
		return Config{}, fmt.Errorf("loading transmission config: %w", err)
	}
	C = cfg
	return cfg, nil
}
