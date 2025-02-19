package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StorageType string `yaml:"storage_type" env-required:"true"`
	Database    `yaml:"db_config"`
	HTTPServer  `yaml:"http_server"`
}

type Database struct {
	DBHost     string `yaml:"db_host"`
	DBUser     string `yaml:"db_user"`
	DBPort     string `yaml:"db_port"`
	DBName     string `yaml:"db_name"`
	DBPassword string `yaml:"db_password"`
	DBssl      string `yaml:"db_ssl"`
}

type HTTPServer struct {
	Addr        string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"10s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatal(err)
	}

	if err := config.validate(); err != nil {
		log.Fatal(err)
	}

	return &config
}

func (c *Config) validate() error {
	if c.StorageType == "postgres" {
		if c.DBHost == "" {
			return errors.New("db_host is required for postgres storageType")
		}
		if c.DBUser == "" {
			return errors.New("db_user is required for postgres storageType")
		}
		if c.DBPort == "" {
			return errors.New("db_port is required for postgres storageType")
		}
		if c.DBName == "" {
			return errors.New("db_name is required for postgres storageType")
		}
		if c.DBPassword == "" {
			return errors.New("db_password is required for postgres storageType")
		}
		if c.DBssl == "" {
			return errors.New("db_ssl is required for postgres storageType")
		}
	}

	return nil
}
