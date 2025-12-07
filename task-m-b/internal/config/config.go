package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	Database   string `yaml:"database" env-required:"true" env:"POSTGRES_URL"`
	HTTPServer `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60"`
}

func MustLoad() *Config {
	// Resolve config path: prefer CONFIG_PATH, fallback to local file
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config path does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err)
		return nil
	}

	// Optional overrides from env for containerized deployments
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.Database = dbURL
	}
	if addr := os.Getenv("HTTP_ADDRESS"); addr != "" {
		cfg.HTTPServer.Address = addr
	} else if port := os.Getenv("PORT"); port != "" {
		cfg.HTTPServer.Address = "0.0.0.0:" + port
	}

	return &cfg
}
