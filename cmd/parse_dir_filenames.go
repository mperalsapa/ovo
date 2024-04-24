package main

import (
	"log"
	"ovo-server/internal/file"
)

var (
	ScanDir = "./testing_media_dir/movies"
)

func main() {
	files := file.ScanPath(ScanDir)

	for _, f := range files {
		// log.Println(f)
		parsed := file.ParseFilename(f)

		log.Printf("Name: %s, Year: %d, MetaProvider: %s, MetaId: %s", parsed.Name, parsed.Year, parsed.MetaProvider, parsed.MetaId)

	}
}
