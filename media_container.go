package main

import "encoding/xml"

type PlexClient struct {
	ServerUrl string
	Token     string
}

type MediaContainer struct {
	XMLName             xml.Name    `xml:"MediaContainer"`
	Text                string      `xml:",chardata"`
	Size                string      `xml:"size,attr"`
	AllowSync           string      `xml:"allowSync,attr"`
	Art                 string      `xml:"art,attr"`
	Identifier          string      `xml:"identifier,attr"`
	LibrarySectionID    string      `xml:"librarySectionID,attr"`
	LibrarySectionTitle string      `xml:"librarySectionTitle,attr"`
	LibrarySectionUUID  string      `xml:"librarySectionUUID,attr"`
	MediaTagPrefix      string      `xml:"mediaTagPrefix,attr"`
	MediaTagVersion     string      `xml:"mediaTagVersion,attr"`
	Thumb               string      `xml:"thumb,attr"`
	Title1              string      `xml:"title1,attr"`
	Title2              string      `xml:"title2,attr"`
	ViewGroup           string      `xml:"viewGroup,attr"`
	ViewMode            string      `xml:"viewMode,attr"`
	DirectoryList       []Directory `xml:"Directory"`
	VideoList           []Video     `xml:"Video"`
}

type Directory struct {
	XMLName      xml.Name   `xml:"Directory"`
	Count        int        `xml:"count,attr"`
	Key          string     `xml:"key,attr"`
	Title        string     `xml:"title,attr"`
	Art          string     `xml:"art,attr"`
	Composite    string     `xml:"composite,attr"`
	Filters      int        `xml:"filters,attr"`
	Refreshing   int        `xml:"refreshing,attr"`
	Thumb        string     `xml:"thumb,attr"`
	Type         string     `xml:"type,attr"`
	Agent        string     `xml:"agent,attr"`
	Scanner      string     `xml:"scanner,attr"`
	Language     string     `xml:"language,attr"`
	Uuid         string     `xml:"uuid,attr"`
	UpdatedAt    string     `xml:"updatedAt,attr"`
	CreatedAt    string     `xml:"createdAt,attr"`
	AllowSync    int        `xml:"allowSync,attr"`
	LocationList []Location `xml:"Location"`
}

type Location struct {
	XMLName xml.Name `xml:"Location"`
	Id      int      `xml:"id,attr"`
	Path    string   `xml:"path,attr"`
}

type Video struct {
	XMLName                xml.Name   `xml:"Video"`
	Text                   string     `xml:",chardata"`
	RatingKey              string     `xml:"ratingKey,attr"`
	Key                    string     `xml:"key,attr"`
	Guid                   string     `xml:"guid,attr"`
	Studio                 string     `xml:"studio,attr"`
	Type                   string     `xml:"type,attr"`
	Title                  string     `xml:"title,attr"`
	TitleSort              string     `xml:"titleSort,attr"`
	ContentRating          string     `xml:"contentRating,attr"`
	Summary                string     `xml:"summary,attr"`
	Rating                 string     `xml:"rating,attr"`
	AudienceRating         string     `xml:"audienceRating,attr"`
	Year                   string     `xml:"year,attr"`
	Tagline                string     `xml:"tagline,attr"`
	Thumb                  string     `xml:"thumb,attr"`
	Art                    string     `xml:"art,attr"`
	Duration               string     `xml:"duration,attr"`
	OriginallyAvailableAt  string     `xml:"originallyAvailableAt,attr"`
	AddedAt                string     `xml:"addedAt,attr"`
	UpdatedAt              string     `xml:"updatedAt,attr"`
	AudienceRatingImage    string     `xml:"audienceRatingImage,attr"`
	HasPremiumExtras       string     `xml:"hasPremiumExtras,attr"`
	HasPremiumPrimaryExtra string     `xml:"hasPremiumPrimaryExtra,attr"`
	RatingImage            string     `xml:"ratingImage,attr"`
	ViewCount              string     `xml:"viewCount,attr"`
	LastViewedAt           string     `xml:"lastViewedAt,attr"`
	Media                  Media      `xml:"Media"`
	GenreList              []Genre    `xml:"Genre"`
	DirectorList           []Director `xml:"Director"`
	WriterList             []Writer   `xml:"Writer"`
	CountryList            []Country  `xml:"Country"`
	RoleList               []Role     `xml:"Role"`
}

type Media struct {
	Text                  string `xml:",chardata"`
	ID                    string `xml:"id,attr"`
	Duration              string `xml:"duration,attr"`
	Bitrate               string `xml:"bitrate,attr"`
	Width                 string `xml:"width,attr"`
	Height                string `xml:"height,attr"`
	AspectRatio           string `xml:"aspectRatio,attr"`
	AudioChannels         string `xml:"audioChannels,attr"`
	AudioCodec            string `xml:"audioCodec,attr"`
	VideoCodec            string `xml:"videoCodec,attr"`
	VideoResolution       string `xml:"videoResolution,attr"`
	Container             string `xml:"container,attr"`
	VideoFrameRate        string `xml:"videoFrameRate,attr"`
	OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
	AudioProfile          string `xml:"audioProfile,attr"`
	Has64bitOffsets       string `xml:"has64bitOffsets,attr"`
	VideoProfile          string `xml:"videoProfile,attr"`
	PartList              []Part `xml:"Part"`
}

type Part struct {
	Text                  string `xml:",chardata"`
	ID                    string `xml:"id,attr"`
	Key                   string `xml:"key,attr"`
	Duration              string `xml:"duration,attr"`
	File                  string `xml:"file,attr"`
	Size                  string `xml:"size,attr"`
	AudioProfile          string `xml:"audioProfile,attr"`
	Container             string `xml:"container,attr"`
	Has64bitOffsets       string `xml:"has64bitOffsets,attr"`
	OptimizedForStreaming string `xml:"optimizedForStreaming,attr"`
	VideoProfile          string `xml:"videoProfile,attr"`
}

type Genre struct {
	XMLName xml.Name `xml:"Genre"`
	Tag     string   `xml:"tag,attr"`
}

type Writer struct {
	XMLName xml.Name `xml:"Writer"`
	Tag     string   `xml:"tag,attr"`
}

type Country struct {
	XMLName xml.Name `xml:"Country"`
	Tag     string   `xml:"tag,attr"`
}

type Role struct {
	XMLName xml.Name `xml:"Role"`
	Tag     string   `xml:"tag,attr"`
}

type Director struct {
	XMLName xml.Name `xml:"Director"`
	Tag     string   `xml:"tag,attr"`
}
