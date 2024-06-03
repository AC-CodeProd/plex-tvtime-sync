package storage

import (
	"plex-tvtime-sync/dto"
	"plex-tvtime-sync/pkg/lib"
	"sync"

	"go.uber.org/fx"
)

type protobufStorage struct {
	logger     lib.Logger
	config     lib.Config
	loadedData *dto.IntMap
	mu         sync.Mutex
}

type ProtobufStorageParams struct {
	fx.In

	Logger lib.Logger
	Config lib.Config
}
