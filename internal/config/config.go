package config

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
	"url-service/internal/constants/storageType"
)

type Config struct {
	Env        string `env:"ENV,required"`
	Database   Database
	HTTPServer HTTPServer
}

type Database struct {
	DBHost     string `env:"DB_HOST"`
	DBUser     string `env:"DB_USER"`
	DBPort     string `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBPassword string `env:"DB_PASSWORD"`
	DBssl      string `env:"DB_SSL"`
}

type HTTPServer struct {
	Addr        string        `env:"HTTP_ADDRESS" envDefault:"0.0.0.0:8080"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT" envDefault:"5s"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"10s"`
}

func MustLoad(stType string) *Config {
	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatal("Failed to parse environment variables: ", err)
	}

	if err := config.validate(stType); err != nil {
		log.Fatal(err)
	}

	return &config
}

func (c *Config) validate(stType string) error {
	if stType == storageType.Postgres {
		if c.Database.DBHost == "" {
			return errors.New("DB_HOST is required for postgres storageType")
		}
		if c.Database.DBUser == "" {
			return errors.New("DB_USER is required for postgres storageType")
		}
		if c.Database.DBPort == "" {
			return errors.New("DB_PORT is required for postgres storageType")
		}
		if c.Database.DBName == "" {
			return errors.New("DB_NAME is required for postgres storageType")
		}
		if c.Database.DBPassword == "" {
			return errors.New("DB_PASSWORD is required for postgres storageType")
		}
		if c.Database.DBssl == "" {
			return errors.New("DB_SSL is required for postgres storageType")
		}
	}
	return nil
}
