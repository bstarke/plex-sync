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
	"net"
	"plex-sync/plex"
	"plex-sync/repository"
	"strconv"
	"strings"
)

// var movie_resolution, movie_format, show_resolution, show_format string
var localHost, remoteHost, localMoviePath, localShowPath, remoteMoviePath, remoteShowPath, thisMachineName string
var sftpl, sftpr *sftpClient
var db *gorm.DB
var filesSql = `SELECT s.host_name AS hostName, v.type AS videoType, vfp.file_path AS filePath
				  FROM video_file_parts vfp
				 INNER JOIN video_files vf on vf.id = vfp.video_file_id
				 INNER JOIN videos v on v.id = vf.video_id
				 INNER JOIN servers s on s.id = v.server_id
				 WHERE NOT EXISTS(SELECT *
									FROM video_file_parts vfps
								   INNER JOIN video_files vfs on vfs.id = vfps.video_file_id
								   INNER JOIN videos vs on vs.id = vfs.video_id
								   INNER JOIN servers ss on ss.id = vs.server_id
								   WHERE ss.id != s.id
									 AND vfps.file_path = vfp.file_path
									 AND vfps.file_size = vfp.file_size)
				ORDER BY s.host_name, v.type, vfp.file_path;`

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
		setup()
		clearDB()
		updateDatabaseForLocal()
		updateDatabaseForRemote()
		copyMissingFiles()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
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

func copyMissingFiles() {
	rows, err := db.Raw(filesSql).Rows()
	if err != nil {
		log.Fatalf("Failed FilesSQL: %v", err)
	}
	defer rows.Close()
	var hostName, videoType, filePath string
	for rows.Next() {
		err = rows.Scan(&hostName, &videoType, &filePath)
		if err != nil {
			log.Printf("Scan Error : %v", err)
		}
		file := missingFiles{
			hostName:  hostName,
			videoType: videoType,
			filePath:  filePath,
		}
		destinationPath := getDestinationPath(file)
		srcClient, dstClient := getSftpClient(file)
		//if localHost == thisMachineName && file.hostName != localHost {
		//	sftpr.Get(destinationPath, file.filePath)
		//} else {
		//	if localHost != thisMachineName {
		//		tempFile, err := os.CreateTemp("/tmp", "plex")
		//		if err != nil {
		//			log.Fatal(err)
		//		}
		//		//defer os.Remove(tempFile.Name())
		//		err = srcClient.Get(tempFile.Name(), file.filePath)
		//		if err != nil {
		//			continue
		//		}
		//		file.filePath = tempFile.Name()
		//	}
		//workDir, _ := os.Getwd()
		//fmt.Printf("Working Directory : %v\n", workDir)
		fmt.Printf("scp %+v:%+v %+v:%+v/\n", srcClient.host, file.filePath, dstClient.host, destinationPath)
		//err := dstClient.Put(file.filePath, destinationPath)
		//if err != nil {
		//	log.Printf("Error putting file : %v", err)
		//}
		//}
	}
}

func getDestinationPath(file missingFiles) (destinationPath string) {
	if file.videoType == "movie" && file.hostName == localHost {
		destinationPath = remoteMoviePath
	} else if file.videoType == "show" && file.hostName == localHost {
		destinationPath = remoteShowPath
	} else if file.videoType == "movie" && file.hostName != localHost {
		destinationPath = localMoviePath
	} else if file.videoType == "show" && file.hostName != localHost {
		destinationPath = localShowPath
	}
	return
}

func getSftpClient(file missingFiles) (srcClient *sftpClient, dstClient *sftpClient) {
	if file.hostName == localHost {
		dstClient = sftpr
		srcClient = sftpl
	} else {
		dstClient = sftpl
		srcClient = sftpr
	}
	return
}

func saveVideo(pvideo plex.Video, hostName string) {
	title := fmt.Sprintf("%v %v %v", pvideo.GrandparentTitle, pvideo.ParentTitle, pvideo.Title)
	title = strings.TrimSpace(title)
	//if strings.HasPrefix(pvideo.GUID, "local://") {
	//	log.Printf("Host %v has %v '%v' with plex guid '%v'", hostName, pvideo.Type, title, pvideo.GUID)
	//}
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

func GetLocalhostFQDN() (fqdn string, err error) {
	ifaces, _ := net.InterfaceAddrs()
	for _, iface := range ifaces {
		var addr net.IP
		addr, _, err = net.ParseCIDR(iface.String())
		if err != nil {
			log.Printf("LookupIP failed: %v", err)
			return
		}
		if !addr.IsLoopback() && !addr.IsLinkLocalMulticast() && !addr.IsLinkLocalUnicast() {
			var hosts []string
			hosts, err = net.LookupAddr(addr.String())
			if err != nil || len(hosts) == 0 {
				continue
			} else {
				fqdn = strings.Trim(hosts[0], ".")
				return
			}
		}
	}
	return
}

func setup() {
	var err error
	db, err = gorm.Open(sqlite.Open("plex.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Open DB Failed: %v", err)
	}
	err = db.AutoMigrate(&repository.Server{}, &repository.Video{}, &repository.VideoFile{}, &repository.VideoFilePart{})
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	thisMachineName, err = GetLocalhostFQDN()
	if err != nil {
		log.Fatalf("unable to get this machines FQDN : %q", err)
	}
	localHost = viper.GetString("local.host")
	remoteHost = viper.GetString("remote.host")
	sftpl, _ = NewConn(viper.GetString("local.host"), viper.GetString("local.sftp.user"), viper.GetInt("local.sftp.port"))
	sftpr, _ = NewConn(viper.GetString("remote.host"), viper.GetString("remote.sftp.user"), viper.GetInt("remote.sftp.port"))
}

type missingFiles struct {
	hostName  string
	videoType string
	filePath  string
}
