package interfaces

import "plex-tvtime-sync/domain/entities"

type ITvTimeRepository interface {
	// Search(showName string) ([]entities.TvTimeSearch, error)
	// GetSeriesDataByName(name string) (*entities.Show, error)
	GetSeasonsDataByName(name string) ([]entities.Season, error)
	GetSeasonsDataBySerieId(serieId int64) ([]entities.Season, error)
	// GetSeasonsData(serieId string) ([]entities.Season, error)
	// GetEpisodesData(serieId string) ([]entities.Episode, error)
	MarkAsWatched(episodeId int64) (bool, error)
}

type ITvTimeUsecase interface {
	// Search(showName string) ([]entities.TvTimeSearch, error)
	// GetSeriesDataByName(name string) (*entities.Show, error)
	GetSeasonsDataByName(name string) ([]entities.Season, error)
	GetSeasonsDataBySerieId(serieId int64) ([]entities.Season, error)
	MarkAsWatched(episodeId int64) (bool, error)
}

type ITvTimeApi interface {
	// Search(showName string) ([]entities.TvTimeSearch, error)
	// GetSeriesDataByName(name string) (*entities.Show, error)
	GetSeasonsDataByName(name string) ([]entities.Season, error)
	GetSeasonsDataBySerieId(serieId int64) ([]entities.Season, error)
	// GetSeasonsData(serieId string) ([]entities.Season, error)
	// GetEpisodesData(serieId string) ([]entities.Episode, error)
	MarkAsWatched(episodeId int64) (bool, error)
}
