package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/pkg/lib"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/fx"
)

type TVTimeApi struct {
	config lib.Config
	logger lib.Logger
}

type TVTimeApiParams struct {
	fx.In

	Config lib.Config
	Logger lib.Logger
}

func (tA *TVTimeApi) SearchShow(showName string) ([]entities.TvTimeSearch, error) {
	const names = "__tvtime_api.go__ : SearchShow"
	getId := func(s *goquery.Selection) map[string]string {
		ret := make(map[string]string)
		if idShow, ok := s.Find("a").Attr("href"); ok {
			re := regexp.MustCompile(`(?P<id>\d+)$`)
			match := re.FindStringSubmatch(idShow)
			for i, name := range re.SubexpNames() {
				if i != 0 && name != "" {
					ret[name] = match[i]
				}
			}
		}
		return ret
	}
	var shows []entities.TvTimeSearch
	mUrl := fmt.Sprintf("https://www.tvtime.com/search?q=%s&limit=20", url.QueryEscape(showName))
	tA.logger.Info(fmt.Sprintf("%s | %s", names, mUrl))

	req, err := http.NewRequest("GET", mUrl, nil)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	req.Header.Set("Host", "www.tvtime.com")
	req.Header.Set("Accept-Language", tA.config.TvTime.AcceptLanguage)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	defer closeFile(resp.Body)

	if resp.StatusCode != 200 {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	doc.Find(".search-item").Each(func(i int, s *goquery.Selection) {
		if filter, ok := s.Find("a i").Attr("class"); ok {
			if strings.Contains(filter, "icon-tvst-genre-miniseeries") {
				result := getId(s)
				if _, ok := result["id"]; ok {
					if id, err := strconv.Atoi(result["id"]); err == nil {
						shows = append(shows, entities.TvTimeSearch{
							ID:        id,
							ShowTitle: s.Find("a strong").Text(),
						})
					}
				}
			}
		}
	})
	return shows, nil
}

func (tA *TVTimeApi) GetShow(serieID int) (*entities.Show, error) {
	const names = "__tvtime_api.go__ : GetShow"
	var (
		infoShow entities.Show
		seasons  []entities.Season
	)
	mUrl := fmt.Sprintf("https://www.tvtime.com/fr/show/%d", (serieID))
	tA.logger.Info(fmt.Sprintf("%s | %s", names, mUrl))

	req, err := http.NewRequest("GET", mUrl, nil)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	req.Header.Set("Host", "www.tvtime.com")
	req.Header.Set("Accept-Language", tA.config.TvTime.AcceptLanguage)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}
	defer closeFile(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	doc.Find("div.seasons div.season-content").Each(func(i int, s *goquery.Selection) {
		var episodes []entities.Episode
		s.Find("ul.episode-list li").Each(func(i int, s *goquery.Selection) {
			if href, ok := s.Find("div.infos div.row a:first-child").Attr("href"); ok {
				if idEpisode, err := strconv.Atoi(strings.Split(href, "/")[5]); err == nil {
					episodes = append(episodes, entities.Episode{
						ID:      idEpisode,
						Name:    strings.TrimSpace(s.Find("div.infos div.row a span.episode-name").Text()),
						AirDate: strings.TrimSpace(s.Find("div.infos div.row a span.episode-air-date").Text()),
						Watched: s.Find("a.watched-btn").HasClass("active"),
					})
				}
			}
		})

		seasons = append(seasons, entities.Season{
			Name:     strings.TrimSpace(s.Find("span[itemprop='name']").Text()),
			Episodes: episodes,
		})
	})

	infoShow = entities.Show{
		ID:       serieID,
		Name:     strings.TrimSpace(doc.Find("div.container-fluid div.heading-info").Children().Text()),
		Overview: strings.TrimSpace(doc.Find("div.show-nav").Children().Find("div.overview").Text()),
		Seasons:  seasons,
	}
	return &infoShow, nil
}

func (tA *TVTimeApi) MarkAsWatched(episodeId int) (bool, error) {
	const names = "__tvtime_api.go__ : MarkAsWatched"
	mUrl := fmt.Sprintf("https://www.tvtime.com/watched_episodes?episode_id=%d", episodeId)
	tA.logger.Info(fmt.Sprintf("%s | %s", names, mUrl))
	req, err := http.NewRequest("PUT", mUrl, nil)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}
	req.Header.Set("Host", "www.tvtime.com")
	req.Header.Set("Cookie", fmt.Sprintf("symfony=%s; tvstRemember=%s", tA.config.TvTime.Token.Symfony, tA.config.TvTime.Token.TvstRemember))
	req.Header.Set("Accept-Language", tA.config.TvTime.AcceptLanguage)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}
	defer closeFile(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return false, err
	}

	// assuming that the response body is a JSON object and "result" is one of its keys
	return string(body) == `{"result":"OK"}`, nil
}

func GetTVTimeApi(taP TVTimeApiParams) TVTimeApi {
	return TVTimeApi{
		config: taP.Config,
		logger: taP.Logger,
	}
}
