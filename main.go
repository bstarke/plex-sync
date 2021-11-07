package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	if xmlBytes, err := getXML("http://nas.home.starkenberg.net:32400/library/sections/1/all?X-Plex-Token=7XzynkzVHNjxz5m_pssP"); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		var result MediaContainer
		xml.Unmarshal(xmlBytes, &result)
		for _, v := range result.Video {
			imdbID := getImdbId(v.Media.Part[0].File)
			if len(imdbID) > 1 {
				sendImdbID(imdbID)
			}
		}
	}
}

func getImdbId(filename string) string {
	part1 := strings.FieldsFunc(filename, func(r rune) bool {
		if r == '{' || r == '}' {
			return true
		}
		return false
	})
	if len(part1) > 1 {
		part2 := strings.FieldsFunc(part1[1], func(r rune) bool {
			if r == '-' {
				return true
			}
			return false
		})
		if len(part2) > 1 {
			return part2[1]
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
	url := fmt.Sprintf("http://localhost:8080/v1/plex/%s", imdbID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		log.Printf("Error posting to API : %s -> %s", imdbID, err)
		return
	}
	if resp.StatusCode != 201 {
		log.Printf("Error posting to API : %s", imdbID)
		log.Printf("Status Code : %d", resp.StatusCode)
	}
}
