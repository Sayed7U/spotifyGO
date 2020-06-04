package gospotify

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"os"
	"strings"
	"sort"
)

// type Tracks struct {
// 	TrackList []Track
// }

type Tracks []Track

type Track struct {
	ArtistName  string `json: "artistName"`
	TrackName string `json: "trackName"`
	MsPlayed int64 `json: "msPlayed"`
	EndTime string `json: "endTime"`
}

func OpenJsonTracks(file string,path string) (Tracks, string) {
	files, _  := filepath.Glob(filepath.Join(path, file) + "*")
	var filesString []string
	var tracks Tracks
	for j:=0; j < len(files); j++ {

		jsonFile, err := os.Open(files[j])

		if err != nil {
			fmt.Println(err)
		}
		s := fmt.Sprintf("Successfully opened %s", files[j])
		fmt.Println(s)
		filesString = append(filesString, s)
		defer jsonFile.Close()

		var tempTracks Tracks

		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &tempTracks)

		tracks = append(tracks, tempTracks...)
	}

	return tracks, strings.Join(filesString,"\n")
}


func (T Tracks) TotalTimePlayed(format string) (float64, string) {
	format = strings.ToLower(format)
	var Time float64
	for i, _ := range T {
		Time += float64(T[i].MsPlayed)
	}
	switch format {
	case "days":
		Time = Time / (1000 * 60 * 60 * 24)
	case "hours":
		Time = Time / (1000 * 60 * 60)

	case "minutes":
		Time = Time / (1000 * 60)

	default:
		Time = Time / 1000
		format = "seconds"
	}
	fmt.Printf("Total time played: %.3f %s. \n", Time, 
	format)
	return Time, format
}

func (T Tracks) AverageTimePlayed(format string) (float64, string) {
	format = strings.ToLower(format)
	var Time float64
	for i, _ := range T {
		Time += float64(T[i].MsPlayed)
	}
	switch format {
	case "days":
		Time = Time / (1000 * 60 * 60 * 24)
	case "hours":
		Time = Time / (1000 * 60 * 60)

	case "minutes":
		Time = Time / (1000 * 60)

	default:
		Time = Time / 1000
	}
	AvgTime := Time/float64(len(T))
	fmt.Printf("The average time played: %.3f %s.\n", 
		AvgTime,format)
	return AvgTime, format
}

func (T Tracks) FindArtistTracksNo(artist string) (int, string) {
	var ret_tracks Tracks
	for i, _ := range T {
		if strings.EqualFold(T[i].ArtistName,artist) {
			ret_tracks = append(ret_tracks, T[i])
		}
	}
	if len(ret_tracks) == 0 {
		fmt.Printf("No tracks found for %s.\n", artist)
		return 0, artist
	}
	fmt.Printf("No. of times %s was played: %v.\n",
		artist, len(ret_tracks))
	retArtist := ret_tracks[0].ArtistName
	return len(ret_tracks), retArtist
}

func (T Tracks) ToString() string {
	var trackStrings []string
	for i := range T {
		s := fmt.Sprintf("%s - %s", T[i].ArtistName, T[i].TrackName)
		trackStrings = append(trackStrings, s)
	}
	return strings.Join(trackStrings, "\n")
}

func (T Tracks) FindArtistTracks(artist string) Tracks{
	var ret_strings []string
	var ret_tracks Tracks
	for i, _ := range T {
		if (strings.EqualFold(T[i].ArtistName,artist) &&
		!stringInSlice(T[i].TrackName,ret_strings)) {
			ret_strings = append(ret_strings, T[i].TrackName)
			ret_tracks = append(ret_tracks, T[i])
		}
	}
	fmt.Println("No. of tracks played by",artist,":",
		len(ret_strings))
	fmt.Printf("Tracks played: %s \n", strings.Join(ret_strings,", "))
	return ret_tracks
}

func (T Tracks) FindTrackName(trackname string) (int, string, string){
	var ret_tracks Tracks
	for i, _ := range T {
		if strings.EqualFold(T[i].TrackName,trackname) {
			ret_tracks = append(ret_tracks, T[i])
		}
	}
	artist := ret_tracks[0].ArtistName
	track := ret_tracks[0].TrackName
	s := fmt.Sprintf("No. of times %s by %s was played: %v",
		track, artist, len(ret_tracks))
	fmt.Println(s)
	return len(ret_tracks),track,artist
}


func (T Tracks) FindArtistPlayed() (string,int) {
	dupfreq := Dup_Count(T)
	// pl := rankByWordCount(dupfreq)
	var plays []int
	var mostPlayed string
	var mostPlayedV int
	prev := 0
	for key, value := range dupfreq {
		if value > prev {
			mostPlayed = key
			mostPlayedV = value
			prev = value
		}
		plays = append(plays,value)
	}
	s := fmt.Sprintf("<b>The most played artist is %s:</b> %v.\n",
			mostPlayed, mostPlayedV)
	fmt.Println(s)
	return mostPlayed, mostPlayedV
}

// func main() {
// 	fmt.Println("Welcome to the spotify data analyser.")
// 	tracks := openJsonTracks("StreamingHistory",
// 		"C:\\Users\\Sayed\\Desktop\\my_spotify_data\\MyData")
// 	tracks.TotalTimePlayed("Days")
// 	tracks.FindArtist("Ariana Grande")
// 	tracks.FindArtist("Tame Impala")
// 	tracks.FindArtist("Justin Bieber")
// 	tracks.FindArtist("Jhené Aiko")
// 	tracks.FindArtist("Drake")
// 	tracks.FindTrackName("Borderline")
// 	tracks.FindTrackName("Bad Day")
// 	tracks.FindTrackName("Lost in Yesterday")
// 	tracks.FindTrackName("Imagination")
// 	tracks.FindTrackName("God is a Woman")
// 	tracks.FindArtistTracks("Future")

// 	tracks.FindArtistPlayed()
// 	tracks.AverageTimePlayed("Seconds")
	
// }

 func Dup_Count(T Tracks) map[string]int {
 	var artists []string

 	for i, _ := range T {
 		artists = append(artists, T[i].ArtistName)
 	}

 	duplicate_frequency := make(map[string]int)
	for _, item := range artists {
		// check if the item/element exist in the duplicate_frequency map
 		_, exist := duplicate_frequency[item]

 		if exist {
 			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
 		} else {
 			duplicate_frequency[item] = 1 // else start counting from 1
 		}
	}
 	//fmt.Println(duplicate_frequency)
 	return duplicate_frequency
 }
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
 func rankByWordCount(wordFrequencies map[string]int) PairList{
  pl := make(PairList, len(wordFrequencies))
  i := 0
  for k, v := range wordFrequencies {
    pl[i] = Pair{k, v}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}

type Pair struct {
  Key string
  Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

