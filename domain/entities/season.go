package entities

type Season struct {
	Name     string
	Number   int
	Episodes []Episode
	SeriesId int64
}
