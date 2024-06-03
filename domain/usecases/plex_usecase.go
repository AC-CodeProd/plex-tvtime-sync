package usecases

import (
	"fmt"
	"plex-tvtime-sync/domain/entities"

	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type PlexUsecaseParams struct {
	fx.In

	Logger         lib.Logger
	PlexRepository interfaces.IPlexRepository
}
type plexUsecase struct {
	logger         lib.Logger
	plexRepository interfaces.IPlexRepository
}

func NewPlexUsecase(puP PlexUsecaseParams) interfaces.IPlexUsecase {
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

func (pU *plexUsecase) DownloadParentThumb(baseUrl, plexToken, parentThumbUrl, dir string) (string, error) {
	return pU.plexRepository.DownloadParentThumb(baseUrl, plexToken, parentThumbUrl, dir)
}

func (pR *plexUsecase) DownloadThumb(baseUrl, plexToken, thumbUrl, dir string) (string, error) {
	return pR.plexRepository.DownloadThumb(baseUrl, plexToken, thumbUrl, dir)
}
