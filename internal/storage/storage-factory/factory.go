package storage_factory

import (
	"fmt"
	"url-service/internal/config"
	"url-service/internal/constants/storageType"
	"url-service/internal/storage"
	"url-service/internal/storage/memory"
	"url-service/internal/storage/postgres"
)

func GetStorage(conf *config.Config) (storage.Storage, error) {
	switch conf.StorageType {
	case storageType.Postgres:
		connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			conf.DBHost,
			conf.DBPort,
			conf.DBUser,
			conf.DBPassword,
			conf.DBName,
			conf.DBssl)

		newStorage, err := postgres.NewPostgres(connectionString)
		if err != nil {
			return nil, err
		}
		return newStorage, nil

	case storageType.Memory:
		newStorage, err := memory.NewMemory()
		if err != nil {
			return nil, err
		}
		return newStorage, nil

	default:
		return nil, fmt.Errorf("unknown storageType type: %s", conf.StorageType)
	}
}
