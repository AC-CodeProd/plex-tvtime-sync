package repositories

import (
	"fmt"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type TvTimeRepositoryParams struct {
	fx.In
	TVTimeApi interfaces.ITvTimeApi
	Logger    lib.Logger
}
type tvTimeRepository struct {
	tvTimeAPI interfaces.ITvTimeApi
	logger    lib.Logger
}

// NewTvTimeRepository initialize users repository
func NewTvTimeRepository(ttrP TvTimeRepositoryParams) interfaces.ITvTimeRepository {
	return &tvTimeRepository{
		tvTimeAPI: ttrP.TVTimeApi,
		logger:    ttrP.Logger,
	}
}

// func (tR *tvTimeRepository) Search(showName string) ([]entities.TvTimeSearch, error) {
// 	const names = "__tvtime_repository.go__ : Search"
// 	shows, err := tR.tvTimeAPI.Search(showName)
// 	if err != nil {
// 		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
// 		return nil, err
// 	}
// 	return shows, nil
// }

// func (tR *tvTimeRepository) GetSeriesDataByName(name string) (*entities.Show, error) {
// 	const names = "__tvtime_repository.go__ : GetSeriesDataByName"
// 	show, err := tR.tvTimeAPI.GetSeriesDataByName(name)
// 	if err != nil {
// 		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
// 		return nil, err
// 	}

// 	return show, nil
// }

func (tR *tvTimeRepository) GetSeasonsDataByName(name string) ([]entities.Season, error) {
	// const names = "__tvtime_repository.go__ : GetSeasonsDataByName"
	seasons, err := tR.tvTimeAPI.GetSeasonsDataByName(name)
	if err != nil {
		return nil, err
	}

	return seasons, nil
}

func (tR *tvTimeRepository) GetSeasonsDataBySerieId(serieId int64) ([]entities.Season, error) {
	// const names = "__tvtime_repository.go__ : GetSeasonsDataBySerieId"
	seasons, err := tR.tvTimeAPI.GetSeasonsDataBySerieId(serieId)
	if err != nil {
		return nil, err
	}

	return seasons, nil
}

func (tR *tvTimeRepository) MarkAsWatched(episodeId int64) (bool, error) {
	const names = "__tvtime_repository.go__ : MarkAsWatched"
	good, err := tR.tvTimeAPI.MarkAsWatched(episodeId)
	if err != nil {
		tR.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}

	return good, nil
}
