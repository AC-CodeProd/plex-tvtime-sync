package tvtime

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"plex-tvtime-sync/domain/entities"
	"plex-tvtime-sync/domain/interfaces"
	"plex-tvtime-sync/pkg/share"
	"slices"
	"strconv"
	"strings"
)

func (tA *TVTimeApi) setRequestHeader(r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("%s %s", tA.config.TvTime.TokenType, tA.config.TvTime.Token))
	r.Header.Set("Accept-Language", tA.config.TvTime.AcceptLanguage)
	r.Header.Set("Country-Code", tA.config.TvTime.CountryCode)
	r.Header.Set("Locale", tA.config.TvTime.Locale)
}

func (tA *TVTimeApi) search(name string) (*TvTimeSearchResponse, error) {

	mUrl := fmt.Sprintf("%s%s%s", tA.config.TvTime.BaseUrl, tA.config.TvTime.SearchUrl, url.QueryEscape(name))

	req, err := http.NewRequest("GET", mUrl, nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Set("Host", "www.tvtime.com")
	tA.setRequestHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer share.CloseFile(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %d %s", resp.StatusCode, resp.Status)
	}
	var data TvTimeSearchResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %d %s", resp.StatusCode, resp.Status)
	}

	return &data, nil
}

func (tA *TVTimeApi) seasons(id int64) (*TvTimeDataSeasonsResponse, error) {

	mUrl := fmt.Sprintf("%s%s", tA.config.TvTime.BaseUrl, strings.Replace(tA.config.TvTime.SeasonsUrl, "{SERIE_ID}", strconv.FormatInt(id, 10), 1))

	req, err := http.NewRequest("GET", mUrl, nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Set("Host", "www.tvtime.com")
	tA.setRequestHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer share.CloseFile(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data TvTimeDataSeasonsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (tA *TVTimeApi) GetSeasonsDataByName(name string) ([]entities.Season, error) {
	const names = "__tvtime_api.go__ : GetSeasonsDataByName"
	searchResponse, err := tA.search(name)
	if err != nil {
		return nil, err
	}
	// fmt.Println(searchResponse)
	if searchResponse.Status != "success" {
		return nil, fmt.Errorf(fmt.Sprintf("%s | %s", names, "searchResponse is not success"))
	}
	var tvTimeSearchDataResponse *TvTimeSearchDataResponse
	// var tvTimeDataSeasonsResponse *TvTimeDataSeasonsResponse

	name = share.ClearString(name)
	for _, data := range searchResponse.Data {
		inTranslate := slices.ContainsFunc(data.Translations, func(s string) bool {
			s = share.ClearString(s)
			return strings.EqualFold(s, name)
		})
		isEqual := strings.EqualFold(share.ClearString(data.Name), name)
		if isEqual || inTranslate {
			tvTimeSearchDataResponse = &data
			break
		}
	}
	if tvTimeSearchDataResponse == nil {
		return nil, fmt.Errorf(fmt.Sprintf("%s | %s has no match in response in tv time", names, name))
	}
	switch tvTimeSearchDataResponse.Type {
	case "series":
		return tA.GetSeasonsDataBySerieId(tvTimeSearchDataResponse.ID)
	default:
		return nil, fmt.Errorf(fmt.Sprintf("%s | the type(%s) has not been implemented", names, tvTimeSearchDataResponse.Type))
	}
}

func (tA *TVTimeApi) GetSeasonsDataBySerieId(serieId int64) ([]entities.Season, error) {
	// const names = "__tvtime_api.go__ : GetSeasonsDataBySerieId"
	var seasons []entities.Season
	var episodes []entities.Episode
	seasonsResponse, err := tA.seasons(serieId)
	if err != nil {
		return nil, err
	}

	seasons = make([]entities.Season, 0)
	for _, season := range seasonsResponse.Seasons {
		episodes = make([]entities.Episode, 0)
		if season.Number == 0 {
			continue
		}
		for _, episode := range season.Episodes {
			episodes = append(episodes, entities.Episode{
				ID:      episode.ID,
				Name:    episode.Name,
				Number:  episode.Number,
				Watched: episode.IsWatched,
			})
		}
		seasons = append(seasons, entities.Season{
			Episodes: episodes,
			Name:     seasonsResponse.Name,
			Number:   season.Number,
			SeriesId: serieId,
		})
	}
	return seasons, nil
}

// func (tA *TVTimeApi) GetSeriesDataByName(name string) (*entities.Show, error) {
// 	const names = "__tvtime_api.go__ : GetSeriesDataByName"
// 	searchResponse, err := tA.search(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// fmt.Println(searchResponse)
// 	if searchResponse.Status != "success" {
// 		return nil, fmt.Errorf(fmt.Sprintf("%s | %s", names, "searchResponse is not success"))
// 	}
// 	var mData *TvTimeSearchDataResponse
// 	name = share.ClearString(name)
// 	for _, data := range searchResponse.Data {
// 		inTranslate := slices.ContainsFunc(data.Translations, func(s string) bool {
// 			s = share.ClearString(s)
// 			return strings.EqualFold(s, name)
// 		})
// 		isEqual := strings.EqualFold(share.ClearString(data.Name), name)
// 		if isEqual || inTranslate {
// 			mData = &data
// 			break
// 		}
// 	}
// 	if mData == nil {
// 		return nil, fmt.Errorf(fmt.Sprintf("%s | %s has no match in response in tv time", names, name))
// 	}
// 	// var (
// 	// 	infoShow entities.Show
// 	// 	seasons  []entities.Season
// 	// )
// 	// mUrl := fmt.Sprintf("%s%s%s", tA.config.TvTime.BaseUrl, tA.config.TvTime.SeriesUrl, uuid)
// 	// tA.logger.Info(fmt.Sprintf("%s | %s", names, mUrl))

// 	// req, err := http.NewRequest("GET", mUrl, nil)
// 	// if err != nil {
// 	// 	tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
// 	// 	return nil, err
// 	// }
// 	// tA.setRequestHeader(req)
// 	// // req.Header.Set("Host", "www.tvtime.com")
// 	// // req.Header.Set("Accept-Language", tA.config.TvTime.AcceptLanguage)

// 	// resp, err := http.DefaultClient.Do(req)
// 	// if err != nil {
// 	// 	tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
// 	// 	return nil, err
// 	// }
// 	// defer share.CloseFile(resp.Body)

// 	// body, err := io.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
// 	// 	return nil, fmt.Errorf("error reading response body: %d %s", resp.StatusCode, resp.Status)
// 	// }
// 	// var data TvTimeDataSeriesResponse
// 	// if err := json.Unmarshal(body, &data); err != nil {
// 	// 	return nil, fmt.Errorf("error decoding JSON: %d %s", resp.StatusCode, resp.Status)
// 	// }
// 	// fmt.Println("__________________________", data)
// 	// // doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	// // if err != nil {
// 	// // 	tA.logger.Error(fmt.Sprintf("%s | %s", names, err))
// 	// // 	return nil, err
// 	// // }

// 	// // doc.Find("div.seasons div.season-content").Each(func(i int, s *goquery.Selection) {
// 	// // 	var episodes []entities.Episode
// 	// // 	s.Find("ul.episode-list li").Each(func(i int, s *goquery.Selection) {
// 	// // 		if href, ok := s.Find("div.infos div.row a:first-child").Attr("href"); ok {
// 	// // 			if idEpisode, err := strconv.Atoi(strings.Split(href, "/")[5]); err == nil {
// 	// // 				episodes = append(episodes, entities.Episode{
// 	// // 					ID:      idEpisode,
// 	// // 					Name:    strings.TrimSpace(s.Find("div.infos div.row a span.episode-name").Text()),
// 	// // 					AirDate: strings.TrimSpace(s.Find("div.infos div.row a span.episode-air-date").Text()),
// 	// // 					Watched: s.Find("a.watched-btn").HasClass("active"),
// 	// // 				})
// 	// // 			}
// 	// // 		}
// 	// // 	})

// 	// // 	seasons = append(seasons, entities.Season{
// 	// // 		Name:     strings.TrimSpace(s.Find("span[itemprop='name']").Text()),
// 	// // 		Episodes: episodes,
// 	// // 	})
// 	// // })

// 	// // infoShow = entities.Show{
// 	// // 	ID:       serieID,
// 	// // 	Name:     strings.TrimSpace(doc.Find("div.container-fluid div.heading-info").Children().Text()),
// 	// // 	Overview: strings.TrimSpace(doc.Find("div.show-nav").Children().Find("div.overview").Text()),
// 	// // 	Seasons:  seasons,
// 	// // }
// 	// infoShow := entities.Show{}
// 	// return &infoShow, nil
// 	return nil, nil
// }

func (tA *TVTimeApi) MarkAsWatched(episodeId int64) (bool, error) {
	// const names = "__tvtime_api.go__ : MarkAsWatched"
	mUrl := fmt.Sprintf("%s%s", tA.config.TvTime.BaseUrl, strings.Replace(tA.config.TvTime.WatchedEpisodesUrl, "{EPISODE_ID}", strconv.FormatInt(episodeId, 10), 1))

	req, err := http.NewRequest("POST", mUrl, nil)
	if err != nil {
		return false, err
	}
	// req.Header.Set("Host", "www.tvtime.com")
	tA.setRequestHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer share.CloseFile(resp.Body)

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %d %s", resp.StatusCode, resp.Status)
	}
	var data TvTimeSearchResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return false, fmt.Errorf("error decoding JSON: %d %s", resp.StatusCode, resp.Status)
	}

	return true, nil
}

func NewTVTimeApi(taP TVTimeApiParams) interfaces.ITvTimeApi {
	return &TVTimeApi{
		config:  taP.Config,
		logger:  taP.Logger,
		helpers: taP.Helpers,
	}
}
