package process

import (
	"fmt"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/usecases"
	"plex-tvtime-sync/pkg/lib"
	"time"

	"go.uber.org/fx"
)

type SyncProcess struct {
	logger        lib.Logger
	config        lib.Config
	jsonStorage   lib.JsonStorage
	tvTimeUsecase usecases.TvTimeUsecase
	plexUsecase   usecases.PlexUsecase
	isRunning     bool
}

type SyncProcessParams struct {
	fx.In

	Logger        lib.Logger
	Config        lib.Config
	JsonStorage   lib.JsonStorage
	TvTimeUsecase usecases.TvTimeUsecase
	PlexUsecase   usecases.PlexUsecase
}

func NewSyncProcess(spP SyncProcessParams) SyncProcess {
	return SyncProcess{
		logger:        spP.Logger,
		config:        spP.Config,
		jsonStorage:   spP.JsonStorage,
		tvTimeUsecase: spP.TvTimeUsecase,
		plexUsecase:   spP.PlexUsecase,
	}
}

func (sH SyncProcess) Run() {
	const names = "__sync_process.go__ : Run"
	lastCheck, err := time.Parse("2006-01-02 15:04:05", sH.config.Plex.InitViewedAt)
	if err != nil {
		sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return
	}

	sH.isRunning = true
	sH.start(lastCheck)
	sH.isRunning = false
	ticker := time.NewTicker(time.Duration(sH.config.Timer) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if !sH.isRunning {
			now := time.Now()
			sH.isRunning = true
			sH.start(lastCheck)
			sH.isRunning = false
			lastCheck = now
		}
	}
}

func (sH SyncProcess) start(lastCheck time.Time) {
	const names = "__sync_process.go__ : start"
	sH.logger.Info(fmt.Sprintf("%s | %s", names, lastCheck.Format("2006-01-02 15:04:05")))
	var viewedAt *int64 = new(int64)
	*viewedAt = lastCheck.Unix()
	plexHistorical, err := sH.plexUsecase.GetLibaryHistory(sH.config.Plex.BaseUrl, sH.config.Plex.Token, "viewedAt:desc", viewedAt, &sH.config.Plex.AccountId)
	if err != nil {
		return
	}
	sH.logger.Info(fmt.Sprintf("%s | There are %d episodes to synchronize.", names, len(plexHistorical)))
	for _, item := range plexHistorical {
		sH.logger.Info(fmt.Sprintf("%s | %s - S%dE%d - %s.", names, item.ShowTitle, item.SeasonNumber, item.EpisodeNumber, item.EpisodeTitle))
		sH.findShow(item)
	}
	sH.logger.Info(fmt.Sprintf("%s | Next check in %d minutes.", names, sH.config.Timer))
}

func (sH SyncProcess) findShow(history entities.PlexHistory) {
	const names = "__sync_process.go__ : findShow"
	sH.logger.Info(fmt.Sprintf("%s | %s", names, "Search for correspondence on TVTime..."))

	if exists, err := sH.jsonStorage.HasLink(history.ID); err != nil {
		sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return
	} else if exists {
		link, _, err := sH.jsonStorage.GetLink(history.ID)
		if err != nil {
			sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
			return
		}

		show, err := sH.tvTimeUsecase.GetShow(link)
		if err != nil {
			sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
			return
		}

		seasonIndex := history.SeasonNumber - 1
		episodeIndex := history.EpisodeNumber - 1

		if seasonIndex >= 0 && seasonIndex < len(show.Seasons) {
			season := show.Seasons[seasonIndex]
			if episodeIndex >= 0 && episodeIndex < len(season.Episodes) {
				episode := season.Episodes[episodeIndex]
				_, err := sH.tvTimeUsecase.MarkAsWatched(episode.ID)
				if err != nil {
					sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
					return
				}

				sH.logger.Info(fmt.Sprintf("%s | %s: marked as seen.", names, episode.Name))
			}
		}
	} else {
		searchs, err := sH.tvTimeUsecase.SearchShow(history.ShowTitle)
		if err != nil {
			return
		}
		for _, search := range searchs {
			show, err := sH.tvTimeUsecase.GetShow(search.ID)
			if err != nil {
				continue
			}
			seasonIndex := history.SeasonNumber - 1
			episodeIndex := history.EpisodeNumber - 1
			if seasonIndex < 0 || seasonIndex >= len(show.Seasons) {
				continue
			}

			season := show.Seasons[seasonIndex]
			if episodeIndex < 0 || episodeIndex >= len(season.Episodes) {
				continue
			}

			episode := season.Episodes[episodeIndex]

			if episode.AirDate != history.Date {
				continue
			}
			// sH.logger.Info(fmt.Sprintf("Correspondance: %s(%d) - %s(%d)", show.Name, show.ID, episode.Name, episode.ID))
			if err := sH.jsonStorage.AddLink(history.ID, show.ID); err != nil {
				sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
				continue
			}

			_, err = sH.tvTimeUsecase.MarkAsWatched(episode.ID)
			if err != nil {
				sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
			}

			sH.logger.Info(fmt.Sprintf("%s | %s: marked as seen.", names, episode.Name))
		}
	}
}
