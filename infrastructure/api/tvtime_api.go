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
	helpers lib.Helpers
	config  lib.Config
	logger  lib.Logger
}

type TVTimeApiParams struct {
	fx.In

	Helpers lib.Helpers
	Config  lib.Config
	Logger  lib.Logger
}

func (tA *TVTimeApi) SearchShow(showName string) ([]entities.TvTimeSearch, error) {
	names, _ := tA.helpers.FuncNameAndFile()
	var shows []entities.TvTimeSearch
	url := fmt.Sprintf("https://www.tvtime.com/search?q=%s&limit=20", url.QueryEscape(showName))
	tA.logger.Info(fmt.Sprintf("%s | %s", names, url))

	req, err := http.NewRequest("GET", url, nil)
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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	re := regexp.MustCompile(`(?P<id>\d+)$`)
	doc.Find(".search-item").Each(func(i int, s *goquery.Selection) {
		filter, _ := s.Find("a i").Attr("class")
		if strings.Contains(filter, "icon-tvst-genre-miniseeries") {
			title := s.Find("a strong").Text()
			idShow, _ := s.Find("a").Attr("href")
			match := re.FindStringSubmatch(idShow)
			result := make(map[string]string)
			for i, name := range re.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			id, err := strconv.Atoi(result["id"])
			if err == nil {
				shows = append(shows, entities.TvTimeSearch{
					ID:        id,
					ShowTitle: title,
				})
			}
		}
	})
	return shows, nil
}

func (tA *TVTimeApi) GetShow(serieID int) (*entities.Show, error) {
	names, _ := tA.helpers.FuncNameAndFile()
	url := fmt.Sprintf("https://www.tvtime.com/fr/show/%d", (serieID))
	tA.logger.Info(fmt.Sprintf("%s | %s", names, url))

	req, err := http.NewRequest("GET", url, nil)
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
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
		return nil, err
	}

	var infoShow entities.Show
	var seasons []entities.Season
	doc.Find("div.seasons div.season-content").Each(func(i int, s *goquery.Selection) {
		var episodes []entities.Episode

		s.Find("ul.episode-list li").Each(func(i int, s *goquery.Selection) {
			linkEpisode := s.Find("div.infos div.row a:first-child")

			href, _ := linkEpisode.Attr("href")
			idEpisode, _ := strconv.Atoi(strings.Split(href, "/")[5])
			nameEpisode := strings.TrimSpace(s.Find("div.infos div.row a span.episode-name").Text())
			airEpisode := strings.TrimSpace(s.Find("div.infos div.row a span.episode-air-date").Text())
			watchedBtn := s.Find("a.watched-btn")
			episodeWatched := watchedBtn.HasClass("active")

			episodes = append(episodes, entities.Episode{
				ID:      idEpisode,
				Name:    nameEpisode,
				AirDate: airEpisode,
				Watched: episodeWatched,
			})
		})

		name := s.Find("span[itemprop='name']").Text()
		seasons = append(seasons, entities.Season{
			Name:     name,
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
	names, _ := tA.helpers.FuncNameAndFile()
	url := fmt.Sprintf("https://www.tvtime.com/watched_episodes?episode_id=%d", episodeId)
	tA.logger.Info(fmt.Sprintf("%s | %s", names, url))
	req, err := http.NewRequest("PUT", url, nil)
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
	defer resp.Body.Close()

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
