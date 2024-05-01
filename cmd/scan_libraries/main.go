package main

import (
	"log"
	"ovo-server/internal/config"
	"ovo-server/internal/database"
	"ovo-server/internal/model"
	"ovo-server/internal/tmdb"
)

func main() {
	config.Init()
	tmdb.Init()
	database.Init()

	libraries := model.GetLibraries()

	for _, library := range libraries {
		log.Printf("Scanning library: %s in %s\n", library.Name, library.Paths)
		library.ScanLibrary()
	}

}
