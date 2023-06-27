package usecases

import (
	"fmt"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/repositories"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type TvTimeUsecase interface {
	SearchShow(showName string) ([]entities.TvTimeSearch, error)
	GetShow(serieID int) (*entities.Show, error)
	MarkAsWatched(episodeId int) (bool, error)
}

type TvTimeUsecaseParams struct {
	fx.In

	Logger           lib.Logger
	TvTimeRepository repositories.TvTimeRepository
}

type tvTimeUsecase struct {
	logger           lib.Logger
	tvTimeRepository repositories.TvTimeRepository
}

func NewTvTimeUsecase(ttuP TvTimeUsecaseParams) TvTimeUsecase {
	return &tvTimeUsecase{
		logger:           ttuP.Logger,
		tvTimeRepository: ttuP.TvTimeRepository,
	}
}

func (tU *tvTimeUsecase) SearchShow(showName string) ([]entities.TvTimeSearch, error) {
	const names = "__tvtime_usecase.go__ : SearchShow"
	tvTimeSearchs, err := tU.tvTimeRepository.SearchShow(showName)
	if err != nil {
		tU.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	return tvTimeSearchs, nil
}

func (tU *tvTimeUsecase) GetShow(serieID int) (*entities.Show, error) {
	const names = "__tvtime_usecase.go__ : GetShow"
	show, err := tU.tvTimeRepository.GetShow(serieID)
	if err != nil {
		tU.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	return show, nil
}

func (tU *tvTimeUsecase) MarkAsWatched(episodeId int) (bool, error) {
	const names = "__tvtime_usecase.go__ : MarkAsWatched"
	good, err := tU.tvTimeRepository.MarkAsWatched(episodeId)
	if err != nil {
		tU.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}

	return good, nil
}
