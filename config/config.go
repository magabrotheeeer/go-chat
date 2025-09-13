package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                     string `yaml:"env"`
	StorageConnectionString string `yaml:"storage_connection_string"`
	HTTPServer              struct {
		AddressHTTP     string `yaml:"addresshttp"`
		TimeoutHTTP     string `yaml:"timeouthttp"`
		IdleTimeOutHTTP string `yaml:"idle_timeouthttp"`
	} `yaml:"http_server"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("file %s - does not exist", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}