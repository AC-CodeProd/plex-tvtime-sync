package entities

type Email struct {
	Recipient            string
	Subject              string
	Body                 string
	TemplateName         string
	SectionSuccessEmails [][]SectionSuccessEmail
	SectionErrorEmails   [][]SectionErrorEmail
}

type SectionErrorEmail struct {
}

type SectionSuccessEmail struct {
	CID          string
	Title        string
	EpisodeTitle string
	Data         string
	Align        string
}
