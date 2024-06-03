package plex

import (
	"encoding/xml"
	"plex-tvtime-sync/pkg/lib"

	"go.uber.org/fx"
)

type PlexApi struct {
	logger  lib.Logger
	helpers lib.Helpers
}

type PlexApiParams struct {
	fx.In

	Logger  lib.Logger
	Helpers lib.Helpers
}

type MediaContainer struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Videos  []Video  `xml:"Video"`
}

type Video struct {
	XMLName               xml.Name `xml:"Video"`
	HistoryKey            string   `xml:"historyKey,attr"`
	Key                   string   `xml:"key,attr"`
	RatingKey             int      `xml:"ratingKey,attr"`
	LibrarySectionID      int64    `xml:"librarySectionID,attr"`
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
	ViewedAt              int64    `xml:"viewedAt,attr"`
	AccountID             int64    `xml:"accountID,attr"`
	DeviceID              int64    `xml:"deviceID,attr"`
}
