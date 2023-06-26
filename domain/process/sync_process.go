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
	helpers       lib.Helpers
	logger        lib.Logger
	config        lib.Config
	jsonStorage   lib.JsonStorage
	tvTimeUsecase usecases.TvTimeUsecase
	plexUsecase   usecases.PlexUsecase
}

type SyncProcessParams struct {
	fx.In

	Helpers       lib.Helpers
	Logger        lib.Logger
	Config        lib.Config
	JsonStorage   lib.JsonStorage
	TvTimeUsecase usecases.TvTimeUsecase
	PlexUsecase   usecases.PlexUsecase
}

func NewSyncProcess(spP SyncProcessParams) SyncProcess {
	return SyncProcess{
		helpers:       spP.Helpers,
		logger:        spP.Logger,
		config:        spP.Config,
		jsonStorage:   spP.JsonStorage,
		tvTimeUsecase: spP.TvTimeUsecase,
		plexUsecase:   spP.PlexUsecase,
	}
}

func (sH SyncProcess) Run() {
	names, _ := sH.helpers.FuncNameAndFile()
	lastCheck, err := time.Parse("2006-01-02 15:04:05", sH.config.Plex.InitViewedAt)
	if err != nil {
		sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
	}
	sH.start(lastCheck)
	ticker := time.NewTicker(time.Duration(sH.config.Timer) * time.Minute)

	for range ticker.C {
		sH.start(lastCheck)
		lastCheck = time.Now()
	}
}

func (sH SyncProcess) start(lastCheck time.Time) {
	names, _ := sH.helpers.FuncNameAndFile()
	sH.logger.Info(fmt.Sprintf("%s | %s", names, lastCheck.Format("2006-01-02 15:04:05")))
	var viewedAt *int64 = new(int64)
	*viewedAt = lastCheck.Unix()
	plexHistorys, _ := sH.plexUsecase.GetLibaryHistory(sH.config.Plex.BaseUrl, sH.config.Plex.Token, "viewedAt:desc", viewedAt, &sH.config.Plex.AccountId)
	sH.logger.Info(fmt.Sprintf("%s | There are %d episodes to synchronize.", names, len(plexHistorys)))
	for _, item := range plexHistorys {
		sH.logger.Info(fmt.Sprintf("%s | %s - S%dE%d - %s.", names, item.ShowTitle, item.SeasonNumber, item.EpisodeNumber, item.EpisodeTitle))
		sH.findShow(item)
	}
	sH.logger.Info(fmt.Sprintf("%s | Next check in %d minutes.", names, sH.config.Timer))
}

func (sH SyncProcess) findShow(history entities.PlexHistory) {
	names, _ := sH.helpers.FuncNameAndFile()
	sH.logger.Info(fmt.Sprintf("%s | %s", names, "Search for correspondence on TVTime..."))

	if exists, err := sH.jsonStorage.HasLink(history.ID); err != nil {
		sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return
	} else if exists {
		link, _, err := sH.jsonStorage.GetLink(history.ID)
		if err != nil {
			sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
		}

		show, err := sH.tvTimeUsecase.GetShow(link)
		if err != nil {
			sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
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
				}

				sH.logger.Info(fmt.Sprintf("%s | %s: marked as seen.", names, episode.Name))
			}
		}
	} else {
		searchs, _ := sH.tvTimeUsecase.SearchShow(history.ShowTitle)
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
			err = sH.jsonStorage.AddLink(history.ID, show.ID)
			if err != nil {
				sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
			}

			_, err = sH.tvTimeUsecase.MarkAsWatched(episode.ID)
			if err != nil {
				sH.logger.Error(fmt.Sprintf("%s | %s", names, err))
			}

			sH.logger.Info(fmt.Sprintf("%s | %s: marked as seen.", names, episode.Name))
		}
	}
}
