package main

import (
	"fmt"
	"log"
	"ovo-server/internal/file"
)

var (
	MoviesDir = "./testing_media_dir/movies"
	ShowsDir  = "./testing_media_dir/shows"
)

func main() {
	files := file.ScanFiles(MoviesDir)

	log.Printf("Found %d files in %s", len(files), MoviesDir)
	for _, f := range files {
		// log.Println(f)
		parsed := file.ParseFilename(f)

		log.Printf("Name: %s, Year: %d, MetaProvider: %s, MetaId: %s", parsed.Name, parsed.Year, parsed.MetaProvider, parsed.MetaID)
	}
	fmt.Println("")
	log.Println("Scanning shows...")
	log.Printf("Found %d directories in %s", len(file.ScanDirectories(ShowsDir)), ShowsDir)
	for _, d := range file.ScanDirectories(ShowsDir) {
		parsed := file.ParseFilename(d)
		log.Printf("\tName: %s, Year: %d, MetaProvider: %s, MetaId: %s", parsed.Name, parsed.Year, parsed.MetaProvider, parsed.MetaID)

		showDir := ShowsDir + "/" + d
		seasons := file.ScanDirectories(showDir)
		log.Printf("\tFound %d directories in %s", len(seasons), showDir)
		for _, s := range seasons {
			seasonNumber, err := file.ParseSeasonDirname(s)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("\t\tSeason: %d", seasonNumber)

			seasonDir := showDir + "/" + s
			episodes := file.ScanFiles(seasonDir)
			log.Printf("\t\tFound %d files in %s", len(episodes), seasonDir)
			for _, e := range episodes {
				parsed, err := file.ParseEpisodeFilename(e)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("\t\t\tEpisode: %d", parsed)
			}
		}
	}
}
