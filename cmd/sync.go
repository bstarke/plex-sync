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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"plex-sync/plex"
	"plex-sync/repository"
	"strconv"
	"strings"
)

// var movie_resolution, movie_format, show_resolution, show_format string
var localHost, remoteHost, localMoviePath, localShowPath, remoteMoviePath, remoteShowPath string
var db *gorm.DB

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
		var err error
		db, err = gorm.Open(sqlite.Open("plex.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Open DB Failed: %v", err)
		}
		err = db.AutoMigrate(&repository.Server{})
		if err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}
		err = db.AutoMigrate(&repository.Video{})
		if err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}
		err = db.AutoMigrate(&repository.VideoFile{})
		if err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}
		err = db.AutoMigrate(&repository.VideoFilePart{})
		if err != nil {
			log.Fatalf("Failed to create tables: %v", err)
		}
		localHost = viper.GetString("local.host")
		remoteHost = viper.GetString("remote.host")
		clearDB()
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
		localHost, viper.GetInt("local.port"))
	pl, _ := plex.New(plexHost, viper.GetString("local.token"))
	dirs, _ := pl.GetLibraries()
	for _, dir := range dirs {
		if dir.Type == "movie" {
			localMoviePath = dir.Location[0].Path
		} else if dir.Type == "show" {
			localShowPath = dir.Location[0].Path
		}
	}
	videos, err := pl.GetAllMovies()
	if err != nil {
		log.Fatalf("Failed to get local movies: %v", err)
	}
	for _, video := range videos {
		saveVideo(video, localHost)
	}
	videos, err = pl.GetAllShows()
	if err != nil {
		log.Fatalf("Failed to get local shows: %v", err)
	}
	for _, video := range videos {
		saveVideo(video, localHost)
	}
}

func updateDatabaseForRemote() {
	plexHost := fmt.Sprintf("%s://%s:%v", viper.Get("remote.protocol"),
		remoteHost, viper.GetInt("remote.port"))
	pr, _ := plex.New(plexHost, viper.GetString("remote.token"))
	dirs, _ := pr.GetLibraries()
	for _, dir := range dirs {
		if dir.Type == "movie" {
			remoteMoviePath = dir.Location[0].Path
		} else if dir.Type == "show" {
			remoteShowPath = dir.Location[0].Path
		}
	}
	videos, err := pr.GetAllMovies()
	if err != nil {
		log.Fatalf("Failed to get remote movies: %v", err)
	}
	for _, video := range videos {
		saveVideo(video, remoteHost)
	}
	videos, err = pr.GetAllShows()
	if err != nil {
		log.Fatalf("Failed to get Remote Shows: %v", err)
	}
	for _, video := range videos {
		saveVideo(video, remoteHost)
	}
}

func fetchMissingMovies() {

}

func pushMissingMovies() {

}

func saveVideo(pvideo plex.Video, hostName string) {
	title := fmt.Sprintf("%v %v %v", pvideo.GrandparentTitle, pvideo.ParentTitle, pvideo.Title)
	title = strings.TrimSpace(title)
	if strings.HasPrefix(pvideo.GUID, "local://") {
		log.Printf("Host %v has %v '%v' with plex guid '%v'", hostName, pvideo.Type, title, pvideo.GUID)
	}
	var server repository.Server
	db.FirstOrCreate(&server, repository.Server{HostName: hostName})
	var video = repository.Video{
		PlexGuid: pvideo.GUID,
		Title:    title,
		Type:     pvideo.Type,
		Year:     pvideo.Year,
	}
	var files []repository.VideoFile
	for _, media := range pvideo.Media {
		var parts []repository.VideoFilePart
		for _, part := range media.Part {
			parts = append(parts, repository.VideoFilePart{
				FilePath: part.File,
				FileSize: part.Size,
			})
		}
		files = append(files, repository.VideoFile{
			Resolution:  getResolution(media.VideoResolution),
			AspectRatio: float32(media.AspectRatio),
			Parts:       parts,
		})
	}
	video.Files = files
	server.Videos = append(server.Videos, video)
	db.Updates(&server)
}

func getResolution(res string) int {
	i, err := strconv.Atoi(res)
	if err == nil {
		return i
	}
	if res == "hd" {
		return 1080
	} else if res == "4k" {
		return 4000
	}
	return 480
}

func clearDB() {
	db.Exec("DELETE FROM video_file_parts")
	db.Exec("DELETE FROM video_files")
	db.Exec("DELETE FROM videos")
}
