package usecases

import (
	"fmt"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/repositories"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type PlexUsecase interface {
	GetLibaryHistory(baseUrl string, plexToken string, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error)
}

type PlexUsecaseParams struct {
	fx.In

	Logger         lib.Logger
	PlexRepository repositories.PlexRepository
}
type plexUsecase struct {
	logger         lib.Logger
	plexRepository repositories.PlexRepository
}

func NewPlexUsecase(puP PlexUsecaseParams) PlexUsecase {
	return &plexUsecase{
		logger:         puP.Logger,
		plexRepository: puP.PlexRepository,
	}
}

func (pU *plexUsecase) GetLibaryHistory(baseUrl string, plexToken string, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error) {
	const names = "__plex_usecase.go__ : GetLibaryHistory"
	historical, err := pU.plexRepository.GetLibaryHistory(baseUrl, plexToken, sort, viewedAt, accountId)
	if err != nil {
		pU.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	return historical, nil
}
