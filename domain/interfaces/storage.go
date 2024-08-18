package interfaces

import (
	"plex-tvtime-sync/dto"
)

type IStorage interface {
	Save(data *dto.IntMap) error
	Load() (*dto.IntMap, error)
	AddSpecificPair(plexId, tvtimeId *int64) error
	GetValue(key int64) (*int64, error)
}

type IStorageUsecase interface {
	SaveIntMap(data map[int]int) error
	LoadIntMap() (map[int]int, error)
	AddSpecificPair(plexId, tvtimeId *int64) error
	GetValue(key int64) (*int64, error)
}
