package usecases

import (
	"fmt"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type TvTimeUsecaseParams struct {
	fx.In

	Logger           lib.Logger
	TvTimeRepository interfaces.ITvTimeRepository
}

type tvTimeUsecase struct {
	logger           lib.Logger
	tvTimeRepository interfaces.ITvTimeRepository
}

func NewTvTimeUsecase(ttuP TvTimeUsecaseParams) interfaces.ITvTimeUsecase {
	return &tvTimeUsecase{
		logger:           ttuP.Logger,
		tvTimeRepository: ttuP.TvTimeRepository,
	}
}

// func (tU *tvTimeUsecase) Search(showName string) ([]entities.TvTimeSearch, error) {
// 	const names = "__tvtime_usecase.go__ : Search"
// 	tvTimeSearchs, err := tU.tvTimeRepository.Search(showName)
// 	if err != nil {
// 		tU.logger.Error(fmt.Sprintf("%s | %s", names, err))
// 		return nil, err
// 	}

// 	return tvTimeSearchs, nil
// }

// func (tU *tvTimeUsecase) GetSeriesDataByName(name string) (*entities.Show, error) {
// 	const names = "__tvtime_usecase.go__ : GetSeriesDataByName"
// 	show, err := tU.tvTimeRepository.GetSeriesDataByName(name)
// 	if err != nil {
// 		tU.logger.Error(fmt.Sprintf("%s | %s", names, err))
// 		return nil, err
// 	}

// 	return show, nil
// }

func (tU *tvTimeUsecase) GetSeasonsDataByName(name string) ([]entities.Season, error) {
	// const names = "__tvtime_usecase.go__ : GetSeasonsDataByName"
	seasons, err := tU.tvTimeRepository.GetSeasonsDataByName(name)
	if err != nil {
		return nil, err
	}

	return seasons, nil
}

func (tU *tvTimeUsecase) GetSeasonsDataBySerieId(serieId int64) ([]entities.Season, error) {
	// const names = "__tvtime_usecase.go__ : GetSeasonsDataByName"
	seasons, err := tU.tvTimeRepository.GetSeasonsDataBySerieId(serieId)
	if err != nil {
		return nil, err
	}

	return seasons, nil
}

func (tU *tvTimeUsecase) MarkAsWatched(episodeId int64) (bool, error) {
	const names = "__tvtime_usecase.go__ : MarkAsWatched"
	good, err := tU.tvTimeRepository.MarkAsWatched(episodeId)
	if err != nil {
		tU.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}

	return good, nil
}
