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
	Helpers   lib.Helpers
	TVTimeApi api.TVTimeApi
	Logger    lib.Logger
}
type tvTimeRepository struct {
	helpers   lib.Helpers
	tvTimeAPI api.TVTimeApi
	logger    lib.Logger
}

// NewTvTimeRepository initialize users repository
func NewTvTimeRepository(ttrP TvTimeRepositoryParams) TvTimeRepository {
	return &tvTimeRepository{
		helpers:   ttrP.Helpers,
		tvTimeAPI: ttrP.TVTimeApi,
		logger:    ttrP.Logger,
	}
}

func (tR *tvTimeRepository) SearchShow(showName string) ([]entities.TvTimeSearch, error) {
	names, _ := tR.helpers.FuncNameAndFile()
	shows, err := tR.tvTimeAPI.SearchShow(showName)
	if err != nil {
		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	return shows, nil
}

func (tR *tvTimeRepository) GetShow(serieID int) (*entities.Show, error) {
	names, _ := tR.helpers.FuncNameAndFile()
	show, err := tR.tvTimeAPI.GetShow(serieID)
	if err != nil {
		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	return show, nil
}

func (tR *tvTimeRepository) MarkAsWatched(episodeId int) (bool, error) {
	names, _ := tR.helpers.FuncNameAndFile()
	good, err := tR.tvTimeAPI.MarkAsWatched(episodeId)
	if err != nil {
		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}

	return good, nil
}
