package main

import (
	"encoding/json"
	"log"
	"ovo-server/internal/file"
	"time"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type ProbeData struct {
	Format Format `json:"format"`
}

type Format struct {
	Filename       string
	NbStreams      int     `json:"nb_streams"`
	NbPrograms     int     `json:"nb_programs"`
	FormatName     string  `json:"format_name"`
	FormatLongName string  `json:"format_long_name"`
	StartTime      float64 `json:"start_time,string"`
	Duration       float64 `json:"duration,string"`
	Size           uint    `json:"size,string"`
	BitRate        uint    `json:"bit_rate,string"`
}

func (p *ProbeData) Duration() time.Duration {
	duration := time.Duration(p.Format.Duration * float64(time.Second))
	return duration
}

func main() {
	probeString, err := ffmpeg_go.Probe("D:/Users/Marc/Documents/REPO/SAPA/DAW2/PROJECTE/FINAL/projecte-final-daw/ovo/testing_media_dir/movies/Big_Buck_Bunny_(2008)_[tmdbid-10378].mp4")
	if err != nil {
		panic(err)
	}

	probeData := ProbeData{}
	json.Unmarshal([]byte(probeString), &probeData)

	log.Printf("ProbeData: %+v", probeData)
	log.Println("Duration: ", probeData.Duration().Seconds())
	log.Println("Duration: ", probeData.Duration())
	minutes := int(probeData.Duration().Minutes())
	seconds := int(probeData.Duration().Seconds()) - minutes*60
	log.Printf("Duration: %dm %ds", minutes, seconds)

	newProbe := file.GetFileMetadata("testing_media_dir/movies/Big_Buck_Bunny_(2008)_[tmdbid-10378].mp4")
	log.Println(newProbe.FancyDuration())

	newProbe = file.GetFileMetadata("\\\\ARCHIUM\\Jellyfin\\Peliculas\\Golpe_en_la_pequeña_China_(1986)_[tmdbid-6978]_Cast.mp4")
	log.Println(newProbe.FancyDuration())
}
