package repositories

import (
	"fmt"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/infrastructure/api"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type PlexRepository interface {
	GetLibaryHistory(baseUrl string, plexToken string, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error)
}

type PlexRepositoryParams struct {
	fx.In

	PlexApi api.PlexApi
	Logger  lib.Logger
}
type plexRepository struct {
	plexApi api.PlexApi
	logger  lib.Logger
}

// NewPlexRepository initialize users repository
func NewPlexRepository(prP PlexRepositoryParams) PlexRepository {
	return &plexRepository{
		plexApi: prP.PlexApi,
		logger:  prP.Logger,
	}
}

func (pR *plexRepository) GetLibaryHistory(baseUrl string, plexToken string, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error) {
	const names = "__plex_repository.go__ : GetLibaryHistory"
	var historical []entities.PlexHistory

	historical, err := pR.plexApi.GetLibaryHistory(baseUrl, plexToken, sort, viewedAt, accountId)
	if err != nil {
		pR.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return historical, err
	}
	return historical, nil
}
