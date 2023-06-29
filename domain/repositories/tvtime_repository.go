package repositories

import (
	"fmt"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/infrastructure/api"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type TvTimeRepository interface {
	SearchShow(showName string) ([]entities.TvTimeSearch, error)
	GetShow(serieID int) (*entities.Show, error)
	MarkAsWatched(episodeId int) (bool, error)
}
type TvTimeRepositoryParams struct {
	fx.In
	TVTimeApi api.TVTimeApi
	Logger    lib.Logger
}
type tvTimeRepository struct {
	tvTimeAPI api.TVTimeApi
	logger    lib.Logger
}

// NewTvTimeRepository initialize users repository
func NewTvTimeRepository(ttrP TvTimeRepositoryParams) TvTimeRepository {
	return &tvTimeRepository{
		tvTimeAPI: ttrP.TVTimeApi,
		logger:    ttrP.Logger,
	}
}

func (tR *tvTimeRepository) SearchShow(showName string) ([]entities.TvTimeSearch, error) {
	const names = "__tvtime_repository.go__ : SearchShow"
	shows, err := tR.tvTimeAPI.SearchShow(showName)
	if err != nil {
		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	return shows, nil
}

func (tR *tvTimeRepository) GetShow(serieID int) (*entities.Show, error) {
	const names = "__tvtime_repository.go__ : GetShow"
	show, err := tR.tvTimeAPI.GetShow(serieID)
	if err != nil {
		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	return show, nil
}

func (tR *tvTimeRepository) MarkAsWatched(episodeId int) (bool, error) {
	const names = "__tvtime_repository.go__ : MarkAsWatched"
	good, err := tR.tvTimeAPI.MarkAsWatched(episodeId)
	if err != nil {
		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}

	return good, nil
}
