package main

import "testing"

func Test_getImdbId(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{"2 Fast 2 Furious", "/srv/media/Movies/2 Fast 2 Furious (2003) {imdb-tt0322259}.mp4", "tt0322259"},
		{"310 to Yuma", "/srv/media/Movies/310 to Yuma (2007) {imdb-tt0381849}.mp4", "tt0381849"},
		{"The 5Th Wave", "/srv/media/Movies/The 5Th Wave (2016) {imdb-tt2304933}.mp4", "tt2304933"},
		{"The 13th Warrior", "/srv/media/Movies/The 13th Warrior (1999) {imdb-tt0120657}.mp4", "tt0120657"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getImdbId(tt.fileName); got != tt.want {
				t.Errorf("getImdbId() = %v, want %v", got, tt.want)
			}
		})
	}
}
