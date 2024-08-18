package interfaces

import (
	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type CMD struct {
	logger         lib.Logger
	storageUseCase interfaces.IStorageUsecase
}

type CMDParams struct {
	fx.In

	Logger lib.Logger

	StorageUseCase interfaces.IStorageUsecase
}

func NewCMD(cP CMDParams) CMD {
	return CMD{
		logger:         cP.Logger,
		storageUseCase: cP.StorageUseCase,
	}
}

func (c CMD) AddSpecificPair(plexId, tvtimeId *int64) error {
	return c.storageUseCase.AddSpecificPair(plexId, tvtimeId)
}
