package plex

import (
	"encoding/json"
	"errors"
)

type Plex struct {
	ServerUrl string
	Token     string
}

func (c Plex) GetAllMovies() (movies []Video, err error) {
	return nil, errors.New("function not implemented")
}

func (c Plex) GetAllShows() (movies []Video, err error) {
	return nil, errors.New("function not implemented")
}

type Response struct {
	MediaContainer MediaContainer `json:"MediaContainer"`
}

type MediaContainer struct {
	Size                int         `json:"size"`
	AllowSync           bool        `json:"allowSync"`
	Art                 string      `json:"art"`
	Identifier          string      `json:"identifier"`
	LibrarySectionID    int         `json:"librarySectionID"`
	LibrarySectionTitle string      `json:"librarySectionTitle"`
	LibrarySectionUUID  string      `json:"librarySectionUUID"`
	MediaTagPrefix      string      `json:"mediaTagPrefix"`
	MediaTagVersion     int         `json:"mediaTagVersion"`
	Thumb               string      `json:"thumb"`
	Title1              string      `json:"title1"`
	Title2              string      `json:"title2"`
	ViewGroup           string      `json:"viewGroup"`
	ViewMode            int         `json:"viewMode"`
	Videos              []Video     `json:"Metadata"`
	Directory           []Directory `json:"Directory"`
}

type Video struct {
	RatingKey             string       `json:"ratingKey"`
	Key                   string       `json:"key"`
	GUID                  string       `json:"guid"`
	Studio                string       `json:"studio"`
	Type                  string       `json:"type"`
	Title                 string       `json:"title"`
	ContentRating         string       `json:"contentRating"`
	Summary               string       `json:"summary"`
	AudienceRating        float64      `json:"audienceRating"`
	Year                  int          `json:"year"`
	Tagline               string       `json:"tagline"`
	Thumb                 string       `json:"thumb"`
	Art                   string       `json:"art"`
	Duration              int          `json:"duration"`
	OriginallyAvailableAt string       `json:"originallyAvailableAt"`
	AddedAt               int          `json:"addedAt"`
	UpdatedAt             int          `json:"updatedAt"`
	AudienceRatingImage   string       `json:"audienceRatingImage"`
	ChapterSource         string       `json:"chapterSource"`
	Media                 []Media      `json:"Media"`
	Genre                 []TaggedData `json:"Genre"`
	Director              []TaggedData `json:"Director"`
	Writer                []TaggedData `json:"Writer"`
	Country               []TaggedData `json:"Country"`
	Role                  []TaggedData `json:"Role"`
	SkipCount             int          `json:"skipCount,omitempty"`
}

type Media struct {
	ID                    int     `json:"id"`
	Duration              int     `json:"duration"`
	Bitrate               int     `json:"bitrate"`
	Width                 int     `json:"width"`
	Height                int     `json:"height"`
	AspectRatio           float64 `json:"aspectRatio"`
	AudioChannels         int     `json:"audioChannels"`
	AudioCodec            string  `json:"audioCodec"`
	VideoCodec            string  `json:"videoCodec"`
	VideoResolution       string  `json:"videoResolution"`
	Container             string  `json:"container"`
	VideoFrameRate        string  `json:"videoFrameRate"`
	OptimizedForStreaming int     `json:"optimizedForStreaming"`
	AudioProfile          string  `json:"audioProfile"`
	Has64BitOffsets       bool    `json:"has64bitOffsets"`
	VideoProfile          string  `json:"videoProfile"`
	Part                  []Part  `json:"Part"`
}

type Part struct {
	ID                    int    `json:"id"`
	Key                   string `json:"key"`
	Duration              int    `json:"duration"`
	File                  string `json:"file"`
	Size                  int    `json:"size"`
	AudioProfile          string `json:"audioProfile"`
	Container             string `json:"container"`
	Has64BitOffsets       bool   `json:"has64bitOffsets"`
	HasThumbnail          string `json:"hasThumbnail"`
	OptimizedForStreaming bool   `json:"optimizedForStreaming"`
	VideoProfile          string `json:"videoProfile"`
}

type TaggedData struct {
	Tag    string      `json:"tag"`
	Filter string      `json:"filter"`
	ID     json.Number `json:"id"`
}

type Directory struct {
	AllowSync        bool       `json:"allowSync"`
	Art              string     `json:"art"`
	Composite        string     `json:"composite"`
	Filters          bool       `json:"filters"`
	Refreshing       bool       `json:"refreshing"`
	Thumb            string     `json:"thumb"`
	Key              string     `json:"key"`
	Type             string     `json:"type"`
	Title            string     `json:"title"`
	Agent            string     `json:"agent"`
	Scanner          string     `json:"scanner"`
	Language         string     `json:"language"`
	UUID             string     `json:"uuid"`
	UpdatedAt        int        `json:"updatedAt"`
	CreatedAt        int        `json:"createdAt"`
	ScannedAt        int        `json:"scannedAt"`
	Content          bool       `json:"content"`
	Directory        bool       `json:"directory"`
	ContentChangedAt int        `json:"contentChangedAt"`
	Hidden           int        `json:"hidden"`
	Location         []Location `json:"Location"`
}

type Location struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
}
