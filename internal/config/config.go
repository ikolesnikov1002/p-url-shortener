package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"environment" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func LoadConfig() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")

	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set.")
	}

	if _, err := os.Stat(cfgPath); err != nil {
		log.Fatalf("Can't load config file %v", cfgPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(cfgPath, &cfg)

	if err != nil {
		log.Fatalf("Can't load config file: %v", err)
	}

	return &cfg
}
