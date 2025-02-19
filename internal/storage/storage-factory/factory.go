package storage_factory

import (
	"fmt"
	"url-service/internal/config"
	"url-service/internal/constants/storageType"
	"url-service/internal/storage"
	"url-service/internal/storage/memory"
	"url-service/internal/storage/postgres"
)

func GetStorage(stType string, conf *config.Config) (storage.Storage, error) {
	switch stType {
	case storageType.Postgres:
		connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			conf.Database.DBHost,
			conf.Database.DBPort,
			conf.Database.DBUser,
			conf.Database.DBPassword,
			conf.Database.DBName,
			conf.Database.DBssl)

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
		return nil, fmt.Errorf("unknown storageType type: %s", stType)
	}
}
