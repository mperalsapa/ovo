package file

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type Metadata struct {
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

func (p *Metadata) Duration() time.Duration {
	duration := time.Duration(p.Format.Duration * float64(time.Second))
	return duration
}

func (p *Metadata) FancyDuration() string {
	var durationString string

	hours := int(p.Duration().Hours())
	if hours > 0 {
		durationString = fmt.Sprintf("%dh ", hours)
	}

	minutes := int(p.Duration().Minutes()) - hours*60
	if minutes > 0 {
		durationString = durationString + fmt.Sprintf("%dm", minutes)
	}

	return durationString
}

func GetFileMetadata(path string) Metadata {
	probeString, err := ffmpeg_go.Probe(path)
	if err != nil {
		log.Println("Error probing file:", err)
		return Metadata{}
	}

	probeData := Metadata{}
	json.Unmarshal([]byte(probeString), &probeData)

	return probeData
}
