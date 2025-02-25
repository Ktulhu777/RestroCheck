package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"local" env-required:"true"`
	Storage
	HTTPServer
}

type HTTPServer struct {
	Address     string        `env:"SERVER_ADDRESS" env-required:"true"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT" env-required:"true"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT" env-required:"true"`
}

type Storage struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	DBName   string `env:"DB_NAME" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASS" env-required:"true"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
