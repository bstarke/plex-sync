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

	"github.com/spf13/cobra"
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Creates JSON Files from Libraries API",
	Long: `Creates JSON Files from Libraries API:
Saves the raw json from the api to a file

Output is saved in the home directory unless the --dir -d flag contains a valid directory`,
	Run: func(cmd *cobra.Command, args []string) {
		if !shows {
			saveMoviesJson()
		}
		if !movies {
			fmt.Println("Shows not implemented yet!")
		}
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().BoolVarP(&movies, "movies", "m", false, "only creates movies.json")
	jsonCmd.Flags().BoolVarP(&shows, "shows", "s", false, "only creates shows.json")
	jsonCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to save json files")
	jsonCmd.MarkFlagsMutuallyExclusive("movies", "shows")
}

func saveMoviesJson() {
	if len(dir) == 0 {
		dir, _ = os.UserHomeDir()
	}
	f, err := os.Create(filepath.Join(dir, "movies.json"))
	if err != nil {
		fmt.Printf("error opening file")
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	plexHost := fmt.Sprintf("%s://%s:%v", viper.Get("local.protocol"),
		viper.Get("local.host"), viper.GetInt("local.port"))
	p, _ := plex2.New(plexHost, viper.GetString("local.token"))
	p.WriteMovieJsonToFile(w)
	fmt.Printf("JSON written to %v\n", f.Name())
}
