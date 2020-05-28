package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"os"
)


type Track struct {
	ArtistName  string `json: "artistName"`
	TrackName string `json: "trackName"`
	MsPlayed int64 `json: "msPlayed"`
	EndTime string `json: "endTime"`
}

func openJsonTracks(file string,path string) []Track {
	files, _  := filepath.Glob(filepath.Join(path, file) + "*")
	fmt.Println(files)
	var tracks []Track
	for j:=0; j < len(files); j++ {

		jsonFile, err := os.Open(files[j])

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Successfully Opened", files[j])
		defer jsonFile.Close()

		var tempTracks []Track

		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &tempTracks)

		tracks = append(tracks, tempTracks...)
	}

	return tracks
}


func TotalTimePlayed(tracks []Track, format string) float64{
	var Time float64
	for i, _ := range tracks {
		Time += float64(tracks[i].MsPlayed)
	}
	switch format {
	case "Days":
		Time = Time / (1000 * 60 * 60 * 24)
	case "Hours":
		Time = Time / (1000 * 60 * 60)

	case "Minutes":
		Time = Time / (1000 * 60)

	default:
		Time = Time / 1000
	}
	return Time
}
func main() {
	fmt.Println("Welcome to the spotify data analyser.")
	tracks := openJsonTracks("StreamingHistory",
		"C:\\Users\\Sayed\\Desktop\\my_spotify_data\\MyData")
	fmt.Println(TotalTimePlayed(tracks,"Days"))
}

