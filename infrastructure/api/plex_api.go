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
	logger lib.Logger
}

type PlexApiParams struct {
	fx.In

	Logger lib.Logger
}

func (pA *PlexApi) GetLibaryHistory(baseUrl string, plexToken string, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error) {
	const names = "__plex_api.go__ : GetLibaryHistory"
	getId := func(video entities.Video) (int, error) {
		split := strings.Split(video.GrandparentKey, "/")
		idStr := split[len(split)-1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			pA.logger.Error(fmt.Sprintf("%s | %s", names, err))
			return 0, err
		}
		return id, nil
	}
	var (
		container  entities.MediaContainer
		historical []entities.PlexHistory
	)

	mUrl := fmt.Sprintf("%s/status/sessions/history/all?X-Plex-Token=%s&sort=%s", baseUrl, plexToken, sort)
	if viewedAt != nil {
		mUrl = fmt.Sprintf("%s&viewedAt>=%d", mUrl, *viewedAt)
	}
	if accountId != nil {
		mUrl = fmt.Sprintf("%s&accountID=%d", mUrl, *accountId)
	}
	pA.logger.Info(fmt.Sprintf("%s | %s", names, mUrl))

	resp, err := http.Get(mUrl)
	if err != nil {
		pA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	defer closeFile(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		pA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	if err := xml.Unmarshal(body, &container); err != nil {
		pA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	for _, video := range container.Videos {
		if len(video.GrandparentKey) == 0 {
			continue
		}
		id, err := getId(video)
		if err != nil {
			return nil, err
		}
		historical = append(historical, entities.PlexHistory{
			ID:            id,
			EpisodeTitle:  video.Title,
			ShowTitle:     video.GrandparentTitle,
			EpisodeNumber: video.Index,
			SeasonNumber:  video.ParentIndex,
			Date:          video.OriginallyAvailableAt,
			ViewedAt:      video.ViewedAt,
		})
	}

	return historical, nil
}

func GetPlexApi(paP PlexApiParams) PlexApi {
	return PlexApi{
		logger: paP.Logger,
	}
}
