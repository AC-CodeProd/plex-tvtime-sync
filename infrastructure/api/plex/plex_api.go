package plex

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/pkg/share"
	"regexp"
	"strconv"
	"strings"
)

func (pA *PlexApi) GetLibaryHistory(baseUrl string, plexToken string, sort string, viewedAt *int64, accountId *int) ([]entities.PlexHistory, error) {
	// const names = "__plex_api.go__ : GetLibaryHistory"
	getId := func(video Video) (int64, error) {
		split := strings.Split(video.GrandparentKey, "/")
		idStr := split[len(split)-1]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	var (
		container  MediaContainer
		historical []entities.PlexHistory
	)

	mUrl := fmt.Sprintf("%s/status/sessions/history/all?X-Plex-Token=%s&sort=%s", baseUrl, plexToken, sort)
	if viewedAt != nil {
		mUrl = fmt.Sprintf("%s&viewedAt>=%d", mUrl, *viewedAt)
	}
	if accountId != nil {
		mUrl = fmt.Sprintf("%s&accountID=%d", mUrl, *accountId)
	}

	resp, err := http.Get(mUrl)
	if err != nil {
		return nil, err
	}
	defer share.CloseFile(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(body, &container); err != nil {
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
			Date:          video.OriginallyAvailableAt,
			EpisodeNumber: video.Index,
			ID:            id,
			SeasonNumber:  video.ParentIndex,
			ShowTitle:     video.GrandparentTitle,
			Thumb:         video.Thumb,
			ParentThumb:   video.ParentThumb,
			ViewedAt:      video.ViewedAt,
		})
	}

	return historical, nil
}

func (pA *PlexApi) DownloadThumb(baseUrl, plexToken, thumbUrl, dir string) (string, error) {
	var re = regexp.MustCompile(`(?m)[a-zA-Z|\W]+`)
	var substitution = "_"
	fileName := fmt.Sprintf("%s.jpg", strings.TrimPrefix(re.ReplaceAllString(thumbUrl, substitution), "_"))
	filePath := fmt.Sprintf("%s/%s", dir, fileName)
	if err := share.EnsureDir(filePath); err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s%s?X-Plex-Token=%s", baseUrl, thumbUrl, plexToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", nil
	}
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	return filePath, nil
}

func (pA *PlexApi) DownloadParentThumb(baseUrl, plexToken, parentThumbUrl, dir string) (string, error) {
	return pA.DownloadThumb(baseUrl, plexToken, parentThumbUrl, dir)
}

func NewPlexApi(paP PlexApiParams) interfaces.IPlexApi {
	return &PlexApi{
		logger:  paP.Logger,
		helpers: paP.Helpers,
	}
}
