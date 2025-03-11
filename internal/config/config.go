package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string        `yaml:"env" env-default:"local"`
	Storage      Storage       `yaml:"storage"`
	HTTPServer   HTTPServer    `yaml:"http_server"`
	Clients      ClientsConfig `yaml:"clients"`
	AppSecret    string        `yaml:"app_secret" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type Storage struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-default:"5432"`
	DBName   string `yaml:"dbname" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-default:"disable"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
	Insecure     bool          `yaml:"insecure"`
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
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
