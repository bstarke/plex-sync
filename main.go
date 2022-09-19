package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	plex2 "plex-sync/plex"
	"strings"
)

func main() {
	run()
}

func setup() {
	viper.SetEnvPrefix("sync")
	viper.SetConfigName("config")           // name of config file (without extension)
	viper.AddConfigPath(".")                // look for config in the working directory
	viper.AddConfigPath("$HOME/.plex-sync") // call multiple times to add many search paths
	err := viper.ReadInConfig()             // Find and read the config file
	if err != nil {                         // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	viper.AutomaticEnv()
	viper.WatchConfig()
}

func run() {
	setup()
	plexHost := fmt.Sprintf("%s://%s:32400", viper.Get("local.protocol"),
		viper.Get("local.host"))
	p, _ := plex2.New(plexHost, viper.GetString("local.token"))
	videos, err := p.GetAllMovies()
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	writeCsvFile(videos, "movies.csv")
	fmt.Printf("Total Movies,%v", len(videos))
	videos, err = p.GetAllShows()
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	writeCsvFile(videos, "shows.csv")
}

func getImdbId(filename string) string {
	fileNameSlices := strings.FieldsFunc(filename, func(r rune) bool {
		if r == '{' || r == '}' {
			return true
		}
		return false
	})
	if len(fileNameSlices) > 1 {
		imdbSlices := strings.FieldsFunc(fileNameSlices[1], func(r rune) bool {
			if r == '-' {
				return true
			}
			return false
		})
		if len(imdbSlices) > 1 {
			return imdbSlices[1]
		}
	}
	return ""
}

func writeCsvFile(videos []plex2.Video, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error opening file")
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v\n", "Title", "Year", "ImdbId", "Resolution", "Format"))
	for index, video := range videos {
		if index >= 0 {
			// CSV Format
			var res, format string
			if video.Media[0].VideoResolution == "1080" {
				res = "Blu-Ray"
			} else {
				res = "DVD"
			}
			if video.Media[0].AspectRatio < 1.5 {
				format = "Fullscreen"
			} else {
				format = "Widescreen"
			}
			var title string
			if video.GrandparentTitle != "" {
				title = fmt.Sprintf("%v - %v - %v", video.GrandparentTitle, video.ParentTitle, video.Title)
			} else {
				title = video.Title
			}
			w.WriteString(fmt.Sprintf("\"%v\",%v,%v,%v,%v\n", title, video.Year, getImdbId(video.Media[0].Part[0].File), res, format))
		}
	}
	w.Flush()
}
