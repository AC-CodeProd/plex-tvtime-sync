package entities

type Email struct {
	Body                 string
	Recipients           []string
	SectionErrorEmails   [][]SectionErrorEmail
	SectionSuccessEmails [][]SectionSuccessEmail
	Subject              string
	TemplateName         string
}

type SectionErrorEmail struct {
	Align         string
	CID           string
	Data          string
	EpisodeNumber string
	EpisodeTitle  string
	Error         string
	PlexID        int
	Title         string
}

type SectionSuccessEmail struct {
	Align         string
	CID           string
	Data          string
	EpisodeNumber string
	EpisodeTitle  string
	PlexID        int
	Title         string
}
