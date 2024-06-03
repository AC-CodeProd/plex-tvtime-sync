package interfaces

import (
	"plex-tvtime-sync/domain/entities"
)

type IPlexRepository interface {
	GetLibaryHistory(baseUrl, plexToken, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error)
	DownloadParentThumb(baseUrl, plexToken, parentThumbUrl, dir string) (string, error)
	DownloadThumb(baseUrl, plexToken, ThumbUrl, dir string) (string, error)
}

type IPlexUsecase interface {
	GetLibaryHistory(baseUrl, plexToken, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error)
	DownloadParentThumb(baseUrl, plexToken, parentThumbUrl, dir string) (string, error)
	DownloadThumb(baseUrl, plexToken, ThumbUrl, dir string) (string, error)
}

type IPlexApi interface {
	GetLibaryHistory(baseUrl, plexToken, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error)
	DownloadParentThumb(baseUrl, plexToken, parentThumbUrl, dir string) (string, error)
	DownloadThumb(baseUrl, plexToken, ThumbUrl, dir string) (string, error)
}
