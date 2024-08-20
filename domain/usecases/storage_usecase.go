package usecases

import (
	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/dto"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

// type StorageUseCase struct {
// 	Storage interfaces.IStorage
// }

type StorageUseCaseParams struct {
	fx.In

	Logger  lib.Logger
	Storage interfaces.IStorage
}

type storageUseCase struct {
	logger  lib.Logger
	storage interfaces.IStorage
}

func NewStorageUseCase(sup StorageUseCaseParams) interfaces.IStorageUsecase {
	return &storageUseCase{
		logger:  sup.Logger,
		storage: sup.Storage,
	}
}

func (su *storageUseCase) Save(data map[int]int) error {
	intMap := &dto.IntMap{Map: make(map[int64]int64)}
	for k, v := range data {
		intMap.Map[int64(k)] = int64(v)
	}
	return su.storage.Save(intMap)
}

func (su *storageUseCase) GetAllSpecificPair() (map[int]int, error) {
	intMap, err := su.storage.Load()
	if err != nil {
		return nil, err
	}
	result := make(map[int]int)
	for k, v := range intMap.Map {
		result[int(k)] = int(v)
	}
	return result, nil
}

func (su *storageUseCase) AddSpecificPair(plexId, tvtimeId *int64) error {
	return su.storage.AddSpecificPair(plexId, tvtimeId)
}

func (su *storageUseCase) GetValue(plexId int64) (*int64, error) {
	return su.storage.GetValue(plexId)
}
