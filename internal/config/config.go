package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"ENV"`
	Server   `env:"SERVER"`
	Database `env:"DATABASE"`
}

type Server struct {
	Address string `env:"ADDR"`
}

type Database struct {
	DatabaseURL     string        `env:"DATABASE_URL"`
	MaxAttempts     int           `env:"MAX_ATTEMPTS"`
	AttemptDuration time.Duration `env:"ATTEMPT_DURATION"`
}

func MustLoad() *Config {
	configPath := os.Getenv("config_path")
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read cfg")
	}

	return &cfg
}
