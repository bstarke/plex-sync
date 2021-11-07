package main

import "encoding/xml"

type MediaContainer struct {
	XMLName             xml.Name `xml:"MediaContainer"`
	Text                string   `xml:",chardata"`
	Size                string   `xml:"size,attr"`
	AllowSync           string   `xml:"allowSync,attr"`
	Art                 string   `xml:"art,attr"`
	Identifier          string   `xml:"identifier,attr"`
	LibrarySectionID    string   `xml:"librarySectionID,attr"`
	LibrarySectionTitle string   `xml:"librarySectionTitle,attr"`
	LibrarySectionUUID  string   `xml:"librarySectionUUID,attr"`
	MediaTagPrefix      string   `xml:"mediaTagPrefix,attr"`
	MediaTagVersion     string   `xml:"mediaTagVersion,attr"`
	Thumb               string   `xml:"thumb,attr"`
	Title1              string   `xml:"title1,attr"`
	Title2              string   `xml:"title2,attr"`
	ViewGroup           string   `xml:"viewGroup,attr"`
	ViewMode            string   `xml:"viewMode,attr"`
	Video               []struct {
		Text                   string `xml:",chardata"`
		RatingKey              string `xml:"ratingKey,attr"`
		Key                    string `xml:"key,attr"`
		Guid                   string `xml:"guid,attr"`
		Studio                 string `xml:"studio,attr"`
		Type                   string `xml:"type,attr"`
		Title                  string `xml:"title,attr"`
		TitleSort              string `xml:"titleSort,attr"`
		ContentRating          string `xml:"contentRating,attr"`
		Summary                string `xml:"summary,attr"`
		Rating                 string `xml:"rating,attr"`
		AudienceRating         string `xml:"audienceRating,attr"`
		Year                   string `xml:"year,attr"`
		Tagline                string `xml:"tagline,attr"`
		Thumb                  string `xml:"thumb,attr"`
		Art                    string `xml:"art,attr"`
		Duration               string `xml:"duration,attr"`
		OriginallyAvailableAt  string `xml:"originallyAvailableAt,attr"`
		AddedAt                string `xml:"addedAt,attr"`
		UpdatedAt              string `xml:"updatedAt,attr"`
		AudienceRatingImage    string `xml:"audienceRatingImage,attr"`
		HasPremiumExtras       string `xml:"hasPremiumExtras,attr"`
		HasPremiumPrimaryExtra string `xml:"hasPremiumPrimaryExtra,attr"`
		RatingImage            string `xml:"ratingImage,attr"`
		ViewCount              string `xml:"viewCount,attr"`
		LastViewedAt           string `xml:"lastViewedAt,attr"`
		Media                  struct {
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
			Part                  []struct {
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
			} `xml:"Part"`
		} `xml:"Media"`
		Genre []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Genre"`
		Director struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Director"`
		Writer []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Writer"`
		Country []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Country"`
		Role []struct {
			Text string `xml:",chardata"`
			Tag  string `xml:"tag,attr"`
		} `xml:"Role"`
	} `xml:"Video"`
}
