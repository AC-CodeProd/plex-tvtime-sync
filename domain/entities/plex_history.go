package entities

type PlexHistory struct {
	Date          string
	EpisodeNumber int
	EpisodeTitle  string
	ID            int64
	SeasonNumber  int
	ShowTitle     string
	Thumb         string
	ParentThumb   string
	ViewedAt      int64
}
