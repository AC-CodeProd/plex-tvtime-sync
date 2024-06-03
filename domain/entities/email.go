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
	Align        string
	EpisodeTitle string
	Error        string
	Title        string
}

type SectionSuccessEmail struct {
	Align        string
	CID          string
	Data         string
	EpisodeTitle string
	Title        string
}
