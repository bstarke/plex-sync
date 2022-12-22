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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"plex-sync/plex"
)

// var movie_resolution, movie_format, show_resolution, show_format string
var localMoviePath, localShowPath, remoteMoviePath, remoteShowPath string

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs the local plex server with the remote",
	Long: `Sync will do a bi-directional sync with the remote

This means that it will pull down all new movies from remote, 
push new local movies to the remote, pull movies with a higher resolution,
push movies if local resolution is higher, pull movies of the preferred format,
push copy of formats not on the remote (assumes remote is what you want as master),
it will then follow the same process for shows`,
	Run: func(cmd *cobra.Command, args []string) {
		updateDatabaseForLocal()
		updateDatabaseForRemote()
		fetchMissingMovies()
		pushMissingMovies()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	//movie_resolution = viper.GetString("preferences.movies.resolution")
	//movie_format = viper.GetString("preferences.movies.format")
	//show_resolution = viper.GetString("preferences.shows.resolution")
	//show_format = viper.GetString("preferences.shows.format")
}

func updateDatabaseForLocal() {
	plexHost := fmt.Sprintf("%s://%s:%v", viper.Get("local.protocol"),
		viper.Get("local.host"), viper.GetInt("local.port"))
	pl, _ := plex.New(plexHost, viper.GetString("local.token"))
	videos, err := pl.GetAllMovies()
	if err != nil {
		log.Fatalf("Failed to get local movies: %v", err)
	}
	for _, video := range videos {
		log.Printf("Local Movie : %v", video.Title)
	}
	fmt.Printf("Local Movies : %v\n", len(videos))
}

func updateDatabaseForRemote() {
	plexHost := fmt.Sprintf("%s://%s:%v", viper.Get("remote.protocol"),
		viper.Get("remote.host"), viper.GetInt("remote.port"))
	pr, _ := plex.New(plexHost, viper.GetString("remote.token"))
	videos, err := pr.GetAllMovies()
	if err != nil {
		log.Fatalf("Failed to get remote movies: %v", err)
	}
	for _, video := range videos {
		log.Printf("Remote Movie : %v", video.Title)
	}
	fmt.Printf("Remote Movies : %v\n", len(videos))
}

func fetchMissingMovies() {

}

func pushMissingMovies() {

}
