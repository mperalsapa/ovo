package main

import (
	"ovo-server/internal/database"
	"ovo-server/internal/model"
)

func main() {

	database.Init()

	libraries := model.GetLibraries()

	for _, library := range libraries {
		library.ScanLibrary()
	}

}
