package tvtime

import (
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type TVTimeApi struct {
	config  lib.Config
	logger  lib.Logger
	helpers lib.Helpers
}

type TVTimeApiParams struct {
	fx.In

	Config  lib.Config
	Logger  lib.Logger
	Helpers lib.Helpers
}

type TvTimeSearchResponse struct {
	Status string                     `json:"status"`
	Data   []TvTimeSearchDataResponse `json:"data"`
}

type TvTimeSearchDataResponse struct {
	UUID string `json:"uuid"`
	ID   int64  `json:"id"`
	Name string `json:"name"`
	// Poster struct {
	// 	Height   int    `json:"height"`
	// 	URL      string `json:"url"`
	// 	Versions struct {
	// 		Big    string `json:"big"`
	// 		Medium string `json:"medium"`
	// 		Small  string `json:"small"`
	// 	} `json:"versions"`
	// 	Width int `json:"width"`
	// } `json:"poster"`
	// Posters []struct {
	// 	Type     string `json:"type"`
	// 	UUID     string `json:"uuid"`
	// 	URL      string `json:"url"`
	// 	ThumbURL string `json:"thumb_url"`
	// 	Width    int    `json:"width"`
	// 	Height   int    `json:"height"`
	// } `json:"posters"`
	// FollowerCount        int      `json:"follower_count"`
	Type         string   `json:"type"`
	Translations []string `json:"translations"`
	// TranslationsWithLang []string `json:"translations_with_lang"`
	// BlocklistStatus      string   `json:"blocklist_status"`
	// ReportedCount        int      `json:"reported_count"`
	// Country              string   `json:"country"`
}

type TvTimeDataSeriesResponse struct {
	Status string `json:"status"`
	Data   struct {
		AirTime           string `json:"air_time"`
		AiredEpisodeCount int    `json:"aired_episode_count"`
		Characters        []struct {
			ActorID   int    `json:"actor_id"`
			ActorName string `json:"actor_name"`
			ActorUUID string `json:"actor_uuid"`
			ID        int    `json:"id"`
			ImageURL  string `json:"image_url"`
			IsDeleted bool   `json:"is_deleted"`
			Name      string `json:"name"`
			Type      string `json:"type"`
			UUID      string `json:"uuid"`
		} `json:"characters"`
		Country   string `json:"country"`
		DayOfWeek string `json:"day_of_week"`
		Fanart    []struct {
			FavoriteCount int    `json:"favorite_count"`
			ID            int    `json:"id"`
			IsDeleted     bool   `json:"is_deleted"`
			Language      string `json:"language"`
			ThumbURL      string `json:"thumb_url"`
			Type          string `json:"type"`
			URL           string `json:"url"`
			UUID          string `json:"uuid"`
		} `json:"fanart"`
		FirstAirDate string `json:"first_air_date"`
		FirstEpisode struct {
			ID   int    `json:"id"`
			UUID string `json:"uuid"`
		} `json:"first_episode"`
		Genres   []string `json:"genres"`
		ID       int      `json:"id"`
		ImdbID   string   `json:"imdb_id"`
		Language string   `json:"language"`
		Name     string   `json:"name"`
		Network  string   `json:"network"`
		Overview string   `json:"overview"`
		Posters  []struct {
			FavoriteCount int    `json:"favorite_count"`
			ID            int    `json:"id"`
			IsDeleted     bool   `json:"is_deleted"`
			Language      string `json:"language"`
			ThumbURL      string `json:"thumb_url"`
			Type          string `json:"type"`
			URL           string `json:"url"`
			UUID          string `json:"uuid"`
		} `json:"posters"`
		Rating  int `json:"rating"`
		Reviews any `json:"reviews"`
		Runtime int `json:"runtime"`
		Seasons []struct {
			ID           int    `json:"id"`
			Number       int    `json:"number"`
			Translations []any  `json:"translations"`
			UUID         string `json:"uuid"`
		} `json:"seasons"`
		Status   string `json:"status"`
		Timezone string `json:"timezone"`
		Trailers []struct {
			Embeddable bool   `json:"embeddable"`
			ID         int    `json:"id"`
			Name       string `json:"name"`
			Runtime    int    `json:"runtime"`
			ThumbURL   string `json:"thumb_url"`
			Type       string `json:"type"`
			URL        string `json:"url"`
			UUID       string `json:"uuid"`
		} `json:"trailers"`
		Translations    []any  `json:"translations"`
		Type            string `json:"type"`
		UUID            string `json:"uuid"`
		UtcAirTime      string `json:"utc_air_time"`
		UtcFirstAirDate string `json:"utc_first_air_date"`
	} `json:"data"`
}

type TvTimeSeasonResponse struct {
	ID       string `json:"id"`
	Number   int    `json:"number"`
	Episodes []struct {
		IsWatched    bool   `json:"is_watched"`
		WatchedCount int    `json:"watched_count"`
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		Season       struct {
			ID           string `json:"id"`
			Number       int    `json:"number"`
			EpisodeCount int    `json:"episode_count"`
		} `json:"season"`
		Number    int  `json:"number"`
		IsSpecial bool `json:"is_special"`
		Show      struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"show"`
	} `json:"episodes"`
}

type TvTimeDataSeasonsResponse struct {
	Seasons          []TvTimeSeasonResponse `json:"seasons"`
	ID               int64                  `json:"id"`
	Name             string                 `json:"name"`
	NeedsLazyLoading bool                   `json:"needs_lazy_loading"`
	AppVersion       int                    `json:"app_version"`
}

type TvTimeDataEpisodesResponse struct {
	Status string `json:"status"`
	Data   []struct {
		AbsoluteNumber    int    `json:"absolute_number"`
		AirDate           string `json:"air_date"`
		AirTime           string `json:"air_time"`
		Articles          []any  `json:"articles"`
		Characters        []any  `json:"characters"`
		HasAbsoluteNumber bool   `json:"has_absolute_number"`
		ID                int    `json:"id"`
		ImdbID            string `json:"imdb_id"`
		IsAired           bool   `json:"is_aired"`
		IsMovie           bool   `json:"is_movie"`
		IsSpecial         bool   `json:"is_special"`
		LinkedMovieID     int    `json:"linked_movie_id"`
		Name              string `json:"name"`
		NextEpisode       struct {
			ID   int    `json:"id"`
			UUID string `json:"uuid"`
		} `json:"next_episode,omitempty"`
		Number     int     `json:"number"`
		Overview   string  `json:"overview"`
		Rating     float64 `json:"rating"`
		Runtime    int     `json:"runtime"`
		Screenshot struct {
			FavoriteCount int    `json:"favorite_count"`
			ID            int    `json:"id"`
			IsDeleted     bool   `json:"is_deleted"`
			Language      string `json:"language"`
			ThumbURL      string `json:"thumb_url"`
			Type          string `json:"type"`
			URL           string `json:"url"`
			UUID          string `json:"uuid"`
		} `json:"screenshot"`
		Season struct {
			ID           int    `json:"id"`
			Number       int    `json:"number"`
			Translations any    `json:"translations"`
			UUID         string `json:"uuid"`
		} `json:"season"`
		SeriesID     int    `json:"series_id"`
		SeriesUUID   string `json:"series_uuid"`
		Timestamp    int    `json:"timestamp"`
		Trailers     []any  `json:"trailers"`
		Translations []any  `json:"translations"`
		Type         string `json:"type"`
		UUID         string `json:"uuid"`
		UtcAirDate   string `json:"utc_air_date"`
		UtcAirTime   string `json:"utc_air_time"`
		WatchOrder   int    `json:"watch_order"`
	} `json:"data"`
}
