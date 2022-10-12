package cmd

/*
Copyright © 2022 Brad Starkenberg brad@starkenberg.net

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	plex2 "plex-sync/plex"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var movies, shows bool
var dir string

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Creates CSV Files from Libraries",
	Long: `Creates CSV Files from Libraries:

File format : "Title", "Year", "ImdbId", "Resolution", "Format"
"Format" in the output file is Fullscreen or Widescreen

Output is saved in the home directory unless the --dir -d flag contains a valid directory`,
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
	csvCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to save csv files")
	csvCmd.MarkFlagsMutuallyExclusive("movies", "shows")
}

func writeCsvFile(videos []plex2.Video, filename string) {
	if len(dir) == 0 {
		dir, _ = os.UserHomeDir()
	}
	f, err := os.Create(filepath.Join(dir, filename))
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
