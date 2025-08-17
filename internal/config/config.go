package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	Database   string `yaml:"database" env-required:"true"`
	HTTPServer `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60"`
}

func MustLoad() *Config {
	defaultConfigPath := "C:/Users/samat/GolandProjects/task-manager/config/local.yaml"
	if err := os.Setenv("CONFIG_PATH", defaultConfigPath); err != nil {
		log.Fatal(err)
		return nil
	}

	if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
		log.Fatal("Config path does not exist")
	}

	var cfg Config

	err := cleanenv.ReadConfig(defaultConfigPath, &cfg)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &cfg
}
