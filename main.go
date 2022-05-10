package main

import (
	"encoding/xml"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
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
	url := fmt.Sprintf("%s://%s:32400/library/sections/1/all?X-Plex-Token=%s", viper.Get("plex.protocol"),
		viper.Get("plex.host"), viper.Get("plex.token"))
	if xmlBytes, err := getXML(url); err != nil {
		log.Fatalf("Failed to get XML: %v", err)
	} else {
		var result MediaContainer
		err := xml.Unmarshal(xmlBytes, &result)
		if err != nil {
			log.Fatalln(err)
		}
		for _, v := range result.Video {
			imdbID := getImdbId(v.Media.Part[0].File)
			if len(imdbID) > 1 {
				fmt.Printf("Sending %s to api\n", imdbID)
				sendImdbID(imdbID)
			}
		}
	}
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

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func sendImdbID(imdbID string) {
	url := fmt.Sprintf("%s://%s/v1/plex/%s", viper.Get("movies.protocol"), viper.Get("movies.host"), imdbID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		log.Printf("Error posting to API : %s -> %s", imdbID, err)
		return
	}
	if resp.StatusCode != 201 {
		log.Printf("Error posting to API : %s", imdbID)
		log.Printf("Status Code : %d", resp.StatusCode)
	} else {
		fmt.Printf("%s sent to api\n", imdbID)
	}
}
