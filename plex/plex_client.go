package plex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

const (
	applicationJson = "application/json"
	libraryPath     = `/library/sections/%s/all`
)

// Plex contains fields that are required to make
// an api call to your plex server
type Plex struct {
	URL              string
	Token            string
	ClientIdentifier string
	Headers          headers
	HTTPClient       http.Client
}

type headers struct {
	Platform               string
	PlatformVersion        string
	Provides               string
	Product                string
	Version                string
	Device                 string
	ContainerSize          string
	ContainerStart         string
	Token                  string
	Accept                 string
	ContentType            string
	ClientIdentifier       string
	TargetClientIdentifier string
}

func New(baseUrl, token string) (*Plex, error) {
	var p Plex
	if baseUrl == "" || token == "" {
		return &p, errors.New("url & Token are Required")
	}
	p.HTTPClient = http.Client{
		Timeout: 3 * time.Second,
	}
	p.Headers = defaultHeaders()
	p.ClientIdentifier = p.Headers.ClientIdentifier
	p.Headers.ClientIdentifier = p.ClientIdentifier
	_, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return &p, err
	}
	p.URL = baseUrl
	p.Token = token
	return &p, nil
}

func (p *Plex) get(query string, h headers) (*http.Response, error) {
	client := p.HTTPClient
	req, reqErr := http.NewRequest("GET", query, nil)
	if reqErr != nil {
		return &http.Response{}, reqErr
	}
	p.addHeaders(req, h)
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return resp, nil
}

func (p *Plex) GetLibraries() (dirs []Directory, err error) {
	query := fmt.Sprintf("%s/library/sections", p.URL)
	resp, err := p.get(query, p.Headers)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return dirs, errors.New(resp.Status)
	}
	var result Response
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err.Error())
		return
	}
	dirs = result.MediaContainer.Directory
	return
}

func (p *Plex) GetAllMovies() (movies []Video, err error) {
	movies, err = p.GetVideos(fmt.Sprintf(libraryPath, "1"))
	return
}

func (p *Plex) GetAllShows() (shows []Video, err error) {
	shows, err = p.GetVideos(fmt.Sprintf(libraryPath, "2"))
	return
}

func (p *Plex) GetVideos(key string) ([]Video, error) {
	if key == "" {
		return []Video{}, errors.New("key is required")
	}
	var results Response
	query := fmt.Sprintf("%s%s", p.URL, key)
	resp, err := p.get(query, p.Headers)
	if err != nil {
		return []Video{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return []Video{}, fmt.Errorf("server error: %v", resp.Status)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return []Video{}, err
	}
	var vids []Video
	if results.MediaContainer.Videos[0].Type != "episode" && results.MediaContainer.Videos[0].Type != "movie" {
		for _, video := range results.MediaContainer.Videos {
			vidlist, _ := p.GetVideos(video.Key)
			vids = append(vids, vidlist...)
		}
	} else {
		vids = results.MediaContainer.Videos
	}
	return vids, nil
}

func defaultHeaders() headers {
	version := "0.0.1"

	return headers{
		Platform:         runtime.GOOS,
		PlatformVersion:  "0.0.0",
		Product:          "Go Plex Client",
		Version:          version,
		Device:           runtime.GOOS + " " + runtime.GOARCH,
		ClientIdentifier: "go-plex-client-v" + version,
		ContainerSize:    "Plex-Container-Size=50",
		ContainerStart:   "X-Plex-Container-Start=0",
		Accept:           applicationJson,
		ContentType:      applicationJson,
	}
}

func (p *Plex) addHeaders(req *http.Request, h headers) {
	req.Header.Add("Accept", h.Accept)
	req.Header.Add("X-Plex-Platform", h.Platform)
	req.Header.Add("X-Plex-Platform-Version", h.PlatformVersion)
	req.Header.Add("X-Plex-Provides", h.Provides)
	req.Header.Add("X-Plex-Client-Identifier", p.ClientIdentifier)
	req.Header.Add("X-Plex-Product", h.Product)
	req.Header.Add("X-Plex-Version", h.Version)
	req.Header.Add("X-Plex-Device", h.Device)
	req.Header.Add("X-Plex-Token", p.Token)
}

type Response struct {
	MediaContainer MediaContainer `json:"MediaContainer"`
}

type LibrarySections struct {
	MediaContainer struct {
		Directory []Directory `json:"Directory"`
	} `json:"MediaContainer"`
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
	GrandparentTitle      string       `json:"grandparentTitle"`
	ParentTitle           string       `json:"parentTitle"`
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
