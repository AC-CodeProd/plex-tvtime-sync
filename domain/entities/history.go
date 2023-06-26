package entities

import "encoding/xml"

type Video struct {
	XMLName               xml.Name `xml:"Video"`
	HistoryKey            string   `xml:"historyKey,attr"`
	Key                   string   `xml:"key,attr"`
	RatingKey             int      `xml:"ratingKey,attr"`
	LibrarySectionID      int      `xml:"librarySectionID,attr"`
	ParentKey             string   `xml:"parentKey,attr"`
	GrandparentKey        string   `xml:"grandparentKey,attr"`
	Title                 string   `xml:"title,attr"`
	GrandparentTitle      string   `xml:"grandparentTitle,attr"`
	Type                  string   `xml:"type,attr"`
	Thumb                 string   `xml:"thumb,attr"`
	ParentThumb           string   `xml:"parentThumb,attr"`
	GrandparentThumb      string   `xml:"grandparentThumb,attr"`
	GrandparentArt        string   `xml:"grandparentArt,attr"`
	Index                 int      `xml:"index,attr"`
	ParentIndex           int      `xml:"parentIndex,attr"`
	OriginallyAvailableAt string   `xml:"originallyAvailableAt,attr"`
	ViewedAt              int      `xml:"viewedAt,attr"`
	AccountID             int      `xml:"accountID,attr"`
	DeviceID              int      `xml:"deviceID,attr"`
}

type PlexHistory struct {
	ID            int
	EpisodeTitle  string
	ShowTitle     string
	EpisodeNumber int
	SeasonNumber  int
	Date          string
	ViewedAt      int
}

type MediaContainer struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Videos  []Video  `xml:"Video"`
}
