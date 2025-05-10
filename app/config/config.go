package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

const (
	Prefix = "IAM"
)

// Config example
type Config struct {
	Server *Server
	DB     *DB
}

type Server struct {
	Addr string `envconfig:"ADDR"`
}

type DB struct {
	Host         string        `envconfig:"HOST"`
	Port         string        `envconfig:"PORT"`
	User         string        `envconfig:"USER"`
	Password     string        `envconfig:"PASSWORD"`
	Name         string        `envconfig:"NAME"`
	SSLMode      string        `envconfig:"SSL_MODE"`
	PingInterval time.Duration `envconfig:"PING_INTERVAL"`
}

func NewFromEnv() *Config {
	c := Config{}
	envconfig.MustProcess(Prefix, &c)

	return &c
}
