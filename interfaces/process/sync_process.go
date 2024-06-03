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
		sH.logger.Info(fmt.Sprintf("%s | %d: %s - S%dE%d - %s.", names, item.ID, item.ShowTitle, item.SeasonNumber, item.EpisodeNumber, item.EpisodeTitle))
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
									// if !episode.Watched {
									if episode.Watched {
										// sH.tvTimeUsecase.MarkAsWatched(episode.ID)
										messages = append(messages, struct {
											Item   entities.PlexHistory
											Err    error
											Status string
										}{
											Item:   item,
											Err:    nil,
											Status: "MarkAsWatched",
										})
										// fmt.Printf("______________________________________ID:%d Episode:%d Saison:%d\n", episode.ID, episode.Number, season.Number)
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
	if len(messages) > 0 {
		var filePath string
		var err error
		// for _, message := range messages {
		// 	var chunks [][]var chunks [][]string
		// 	if len(sectionSuccessEmails) == 2 {
		// 		totalSectionSuccessEmails = append(totalSectionSuccessEmails, sectionSuccessEmails)
		// 		sectionSuccessEmails = make([]*entities.SectionSuccessEmail, 2)
		// 	}
		// 	sectionSuccessEmail := entities.SectionSuccessEmail{
		// 		CID:   imageBase64,
		// 		Title: message.Item.ShowTitle,
		// 	}
		// 	sectionSuccessEmails = append(sectionSuccessEmails, &sectionSuccessEmail)
		// 	// body.WriteString(fmt.Sprintf("%s\r\n", fmt.Sprintf(`<img src="cid:%s">`, imageBase64)))
		// }
		chunkSize := 2
		cpt := 0
		for i := 0; i < len(messages); i += chunkSize {
			end := i + chunkSize
			if end > len(messages) {
				end = len(messages)
			}
			tmpMessages := messages[i:end]
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
				cpt++
				cid := fmt.Sprintf("imageCID%d", cpt)
				sectionSuccessEmail := entities.SectionSuccessEmail{
					CID:          cid,
					Data:         imageBase64,
					Title:        message.Item.ShowTitle,
					EpisodeTitle: message.Item.EpisodeTitle,
					// Align: "left",
					// Align: "right",
				}
				m := cpt % 2
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
		// sH.emailUsecase.SendEmailWithTemplate(&entities.Email{
		// 	Subject:              subject,
		// 	Recipient:            sH.config.Mailer.Recipient,
		// 	SectionSuccessEmails: totalSectionSuccessEmails,
		// }, map[string]interface{}{
		// 	"Subject":              subject,
		// 	"SectionSuccessEmails": totalSectionSuccessEmails,
		// 	"Username":             "Alain",
		// 	// "SectionSuccessEmails": totalSectionSuccessEmails,
		// })
	}
	// for _, v := range totalSectionSuccessEmails {
	// 	for _, w := range v {

	// 		fmt.Println(w.Align)
	// 	}
	// }
	// fmt.Println("Iciiiiiiiiiiiiiiiiiiiiii", len(totalSectionSuccessEmails), totalSectionSuccessEmails)
	sH.logger.Info(fmt.Sprintf("%s | Next check in %d minutes.", names, sH.config.Timer))
}
