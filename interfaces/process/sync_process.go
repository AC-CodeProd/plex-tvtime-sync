package process

import (
	"encoding/base64"
	"fmt"
	"os"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/pkg/lib"
	"regexp"
	"time"

	"go.uber.org/fx"
)

type SyncProcess struct {
	logger         lib.Logger
	config         lib.Config
	tvTimeUsecase  interfaces.ITvTimeUsecase
	plexUsecase    interfaces.IPlexUsecase
	emailUsecase   interfaces.IEmailUsecase
	storageUseCase interfaces.IStorageUsecase
	isRunning      bool
}

type SyncProcessParams struct {
	fx.In

	Logger         lib.Logger
	Config         lib.Config
	TvTimeUsecase  interfaces.ITvTimeUsecase
	PlexUsecase    interfaces.IPlexUsecase
	EmailUsecase   interfaces.IEmailUsecase
	StorageUseCase interfaces.IStorageUsecase
}

func NewSyncProcess(spP SyncProcessParams) SyncProcess {
	return SyncProcess{
		logger:         spP.Logger,
		config:         spP.Config,
		tvTimeUsecase:  spP.TvTimeUsecase,
		plexUsecase:    spP.PlexUsecase,
		emailUsecase:   spP.EmailUsecase,
		storageUseCase: spP.StorageUseCase,
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
	plexHistorical, err := sH.plexUsecase.GetLibaryHistory(sH.config.Plex.BaseUrl, sH.config.Plex.Token, sH.config.Plex.ViewedAtSort, viewedAt, &sH.config.Plex.AccountId)
	if err != nil {
		return
	}
	sH.logger.Info(fmt.Sprintf("%s | There are %d episodes to synchronize.", names, len(plexHistorical)))
	var messages = make([]struct {
		Item   entities.PlexHistory
		Err    error
		Status string
	}, 0)
	for _, item := range plexHistorical {
		sH.logger.Debug(fmt.Sprintf("%s | %d: %s - S%dE%d - %s.", names, item.ID, item.ShowTitle, item.SeasonNumber, item.EpisodeNumber, item.EpisodeTitle))
		pattern := `plex://season/([a-z0-9]+)`
		_ = pattern

		if value, err := sH.storageUseCase.GetValue(int64(item.ID)); err != nil {
			messages = append(messages, struct {
				Item   entities.PlexHistory
				Err    error
				Status string
			}{
				Item:   item,
				Err:    err,
				Status: "Error",
			})
			continue
		} else {
			regex, err := regexp.Compile(pattern)
			if err != nil {
				messages = append(messages, struct {
					Item   entities.PlexHistory
					Err    error
					Status string
				}{
					Item:   item,
					Err:    err,
					Status: "Error",
				})
				continue
			}
			matches := regex.FindStringSubmatch(item.ShowTitle)
			// if len(matches) > 0 && value == nil || item.ID == 30013 {
			if len(matches) > 0 && value == nil {
				messages = append(messages, struct {
					Item   entities.PlexHistory
					Err    error
					Status string
				}{
					Item:   item,
					Err:    fmt.Errorf("error %s", "no matches"),
					Status: "Error",
				})
			} else {
				var seasons []entities.Season
				var err error
				if value != nil {
					seasons, err = sH.tvTimeUsecase.GetSeasonsDataBySerieId(*value)
				} else {
					seasons, err = sH.tvTimeUsecase.GetSeasonsDataByName(item.ShowTitle)
				}
				if err != nil {
					messages = append(messages, struct {
						Item   entities.PlexHistory
						Err    error
						Status string
					}{
						Item:   item,
						Err:    err,
						Status: "Error",
					})
					continue
				}

				if len(seasons) < item.SeasonNumber {
					messages = append(messages, struct {
						Item   entities.PlexHistory
						Err    error
						Status string
					}{
						Item:   item,
						Err:    fmt.Errorf("error %s", "the seasons don't match"),
						Status: "Error",
					})
					continue
				} else {
					wasFound := false
					for _, season := range seasons {
						if season.Number == item.SeasonNumber {
							for _, episode := range season.Episodes {
								if episode.Number == item.EpisodeNumber {
									if value == nil {
										if err := sH.storageUseCase.AddSpecificPair(&item.ID, &season.SeriesId); err != nil {
											break
										}
									}
									if !episode.Watched {
										sH.tvTimeUsecase.MarkAsWatched(episode.ID)
										messages = append(messages, struct {
											Item   entities.PlexHistory
											Err    error
											Status string
										}{
											Item:   item,
											Err:    nil,
											Status: "MarkAsWatched",
										})
									}
									wasFound = true
									break
								}
							}
							break
						}
					}
					if !wasFound {
						messages = append(messages, struct {
							Item   entities.PlexHistory
							Err    error
							Status string
						}{
							Item:   item,
							Err:    fmt.Errorf("error %s", "was not found"),
							Status: "Error",
						})
					}
				}
			}
		}
	}
	var totalSectionSuccessEmails [][]entities.SectionSuccessEmail
	var totalSectionErrorEmails [][]entities.SectionErrorEmail
	if len(messages) > 0 {
		chunkSizeSuccess := 2
		chunkSizeError := 2
		var filePath string
		var err error
		cptSuccess := 0
		cptError := 0
		var messagesSuccess = make([]struct {
			Item   entities.PlexHistory
			Err    error
			Status string
		}, 0)
		var messagesError = make([]struct {
			Item   entities.PlexHistory
			Err    error
			Status string
		}, 0)
		for _, message := range messages {
			if message.Status == "Error" {
				messagesError = append(messagesError, message)
			} else {
				messagesSuccess = append(messagesSuccess, message)
			}
		}
		for i := 0; i < len(messagesError); i += chunkSizeError {
			end := i + chunkSizeError
			if end > len(messagesError) {
				end = len(messagesError)
			}
			tmpMessages := messagesError[i:end]
			var sectionErrorEmails = make([]entities.SectionErrorEmail, 0)
			for _, message := range tmpMessages {
				filePath, err = sH.plexUsecase.DownloadParentThumb(sH.config.Plex.BaseUrl, sH.config.Plex.Token, message.Item.ParentThumb, "./tmp/thumb")
				if err != nil {
					fmt.Println(err)
				}
				if filePath == "" {
					filePath, err = sH.plexUsecase.DownloadThumb(sH.config.Plex.BaseUrl, sH.config.Plex.Token, message.Item.Thumb, "./tmp/thumb")
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
				imageBytes, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
				cptError++
				cid := fmt.Sprintf("imageCIDError%d", cptError)
				sectionErrorEmail := entities.SectionErrorEmail{
					CID:           cid,
					Data:          imageBase64,
					EpisodeNumber: fmt.Sprintf("S%dE%d", message.Item.SeasonNumber, message.Item.EpisodeNumber),
					EpisodeTitle:  message.Item.EpisodeTitle,
					Error:         message.Err.Error(),
					PlexID:        int(message.Item.ID),
					Title:         message.Item.ShowTitle,
				}
				m := cptError % 2
				if m != 0 {
					sectionErrorEmail.Align = "left"
				} else {
					sectionErrorEmail.Align = "right"
				}
				sectionErrorEmails = append(sectionErrorEmails, sectionErrorEmail)
			}
			totalSectionErrorEmails = append(totalSectionErrorEmails, sectionErrorEmails)
		}

		for i := 0; i < len(messagesSuccess); i += chunkSizeSuccess {
			end := i + chunkSizeSuccess
			if end > len(messagesSuccess) {
				end = len(messagesSuccess)
			}
			tmpMessages := messagesSuccess[i:end]
			var sectionSuccessEmails = make([]entities.SectionSuccessEmail, 0)
			for _, message := range tmpMessages {
				filePath, err = sH.plexUsecase.DownloadParentThumb(sH.config.Plex.BaseUrl, sH.config.Plex.Token, message.Item.ParentThumb, "./tmp/thumb")
				if err != nil {
					fmt.Println(err)
				}
				if filePath == "" {
					filePath, err = sH.plexUsecase.DownloadThumb(sH.config.Plex.BaseUrl, sH.config.Plex.Token, message.Item.Thumb, "./tmp/thumb")
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
				imageBytes, err := os.ReadFile(filePath)
				if err != nil {
					fmt.Println(err)
					continue
				}
				imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
				cptSuccess++
				cid := fmt.Sprintf("imageCIDSuccess%d", cptSuccess)
				sectionSuccessEmail := entities.SectionSuccessEmail{
					CID:           cid,
					Data:          imageBase64,
					EpisodeNumber: fmt.Sprintf("S%dE%d ", message.Item.SeasonNumber, message.Item.EpisodeNumber),
					EpisodeTitle:  message.Item.EpisodeTitle,
					PlexID:        int(message.Item.ID),
					Title:         message.Item.ShowTitle,
				}
				m := cptSuccess % 2
				if m != 0 {
					sectionSuccessEmail.Align = "left"
				} else {
					sectionSuccessEmail.Align = "right"
				}
				sectionSuccessEmails = append(sectionSuccessEmails, sectionSuccessEmail)
			}
			totalSectionSuccessEmails = append(totalSectionSuccessEmails, sectionSuccessEmails)
		}
		subject := "Plex TVTime Sync: La synchronisation est terminée"
		_ = subject
		email := entities.Email{
			Subject:              subject,
			Recipients:           sH.config.Mailer.Recipients,
			SectionSuccessEmails: totalSectionSuccessEmails,
			SectionErrorEmails:   totalSectionErrorEmails,
		}
		// fmt.Println("icii", totalSectionErrorEmails)
		for _, v := range totalSectionErrorEmails {
			// fmt.Println(v)
			for _, d := range v {
				fmt.Println("PlexID", d.Title, d.PlexID)
			}
		}
		// _ = email
		sH.emailUsecase.SendEmailWithTemplate(&email, map[string]interface{}{
			"Subject":              email.Subject,
			"SectionSuccessEmails": email.SectionSuccessEmails,
			"SectionErrorEmails":   email.SectionErrorEmails,
		})
	}

	sH.logger.Info(fmt.Sprintf("%s | Next check in %d minutes.", names, sH.config.Timer))
}
