package api

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/pkg/lib"
	"strconv"
	"strings"

	"go.uber.org/fx"
)

type PlexApi struct {
	helpers lib.Helpers
	logger  lib.Logger
}

type PlexApiParams struct {
	fx.In

	Helpers lib.Helpers
	Logger  lib.Logger
}

func (pA *PlexApi) GetLibaryHistory(baseUrl string, plexToken string, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error) {
	names, _ := pA.helpers.FuncNameAndFile()
	url := fmt.Sprintf("%s/status/sessions/history/all?X-Plex-Token=%s&sort=%s", baseUrl, plexToken, sort)
	if viewedAt != nil {
		url = fmt.Sprintf("%s&viewedAt>=%d", url, *viewedAt)
	}
	if accountId != nil {
		url = fmt.Sprintf("%s&accountID=%d", url, *accountId)
	}
	pA.logger.Info(fmt.Sprintf("%s | %s", names, url))
	resp, err := http.Get(url)
	if err != nil {
		pA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		pA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	var container entities.MediaContainer
	xml.Unmarshal([]byte(body), &container)
	var historys []entities.PlexHistory
	for _, video := range container.Videos {
		if video.GrandparentKey != "" {
			idStr := strings.Split(video.GrandparentKey, "/")[len(strings.Split(video.GrandparentKey, "/"))-1]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				pA.logger.Error(fmt.Sprintf("%s | %s", names, err))
				return nil, err
			}
			historys = append(historys, entities.PlexHistory{
				ID:            id,
				EpisodeTitle:  video.Title,
				ShowTitle:     video.GrandparentTitle,
				EpisodeNumber: video.Index,
				SeasonNumber:  video.ParentIndex,
				Date:          video.OriginallyAvailableAt,
				ViewedAt:      video.ViewedAt,
			})
		}
	}

	return historys, nil
}

func GetPlexApi(paP PlexApiParams) PlexApi {
	return PlexApi{
		logger: paP.Logger,
	}
}
