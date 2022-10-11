package cmd

/*
Copyright © 2022 Brad Starkenberg brad@starkenberg.net
*/

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	plex2 "plex-sync/plex"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var movies, shows bool

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		plexHost := fmt.Sprintf("%s://%s:%v", viper.Get("local.protocol"),
			viper.Get("local.host"), viper.GetInt("local.port"))
		p, _ := plex2.New(plexHost, viper.GetString("local.token"))
		if !shows {
			videos, err := p.GetAllMovies()
			if err != nil {
				fmt.Printf("error: %v", err)
			}
			writeCsvFile(videos, "movies.csv")
			fmt.Printf("Total Movies : %v\n", len(videos))
		}
		if !movies {
			videos, err := p.GetAllShows()
			if err != nil {
				fmt.Printf("error: %v", err)
			}
			writeCsvFile(videos, "shows.csv")
			fmt.Printf("Total Shows  : %v\n", len(videos))
		}
	},
}

func init() {
	rootCmd.AddCommand(csvCmd)
	csvCmd.Flags().BoolVarP(&movies, "movies", "m", false, "only creates movies.csv")
	csvCmd.Flags().BoolVarP(&shows, "shows", "s", false, "only creates shows.csv")
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
			var res int
			var format string
			if video.Media[0].VideoResolution == "sd" {
				res = 480
			} else {
				res, _ = strconv.Atoi(video.Media[0].VideoResolution)
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
