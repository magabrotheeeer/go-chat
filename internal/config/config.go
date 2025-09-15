package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Volumes  string         `yaml:"-"` // из env
	App      AppConfig
}

type AppConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Name       string `yaml:"name"`
	User       string `yaml:"-"` // из env
	Password   string `yaml:"-"` // из env
	Connection string `yaml:"-"` // из env
	Data       string `yaml:"-"` // из env
	SSLMode    string `yaml:"ssl_mode"`
}

type ServerConfig struct {
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

// MustLoad функция для загрузки конфиг файла
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("file %s - does not exist", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	cfg.Database.User = os.Getenv("POSTGRES_USER")
	if cfg.Database.User == "" {
		log.Fatal("POSTGRES_USER is not set")
	}
	cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	if cfg.Database.Password == "" {
		log.Fatal("POSTGRES_PASSWORD is not set")
	}
	cfg.Database.Connection = os.Getenv("POSTGRES_CONNECTION")
	if cfg.Database.Connection == "" {
		log.Fatal("POSTGRES_CONNECTION is not set")
	}
	cfg.Database.Data = os.Getenv("POSTGRES_DATA")
	if cfg.Database.Data == "" {
		log.Fatal("POSTGRES_DATA is not set")
	}
	cfg.Volumes = os.Getenv("CONFIG_PATH_VOLUMES")
	if cfg.Volumes == "" {
		log.Fatal("CONFIG_PATH_VOLUMES is not set")
	}
	cfg.App.Port = os.Getenv("CHAT_PORT")
	if cfg.App.Port == "" {
		log.Fatal("CHAT_PORT is not set")
	}
	return &cfg
}
