package main

import (
	"fmt"
	"ovo-server/internal/config"
	"ovo-server/internal/database"
	"ovo-server/internal/model"
	"ovo-server/internal/tmdb"
)

func main() {
	config.Init()
	tmdb.Init()
	database.Init()
	model.Init()

	// database.GetDB().Logger.LogMode(logger.Silent)

	library := model.Library{
		Type:  "movie",
		Name:  "Movies",
		Paths: []string{"/movies"},
	}
	library.ID = 1

	database.GetDB().Save(&library)

	item := model.Item{
		LibraryID:     library.ID,
		Title:         "Vengadores End Game",
		OriginalTitle: "Vengadores End Game",
		FilePath:      "/movies/Vengadores End Game (1994).mkv",
	}
	item.ID = 2
	database.GetDB().Save(&item)

	person := model.Person{
		Name:         "Robert Downey",
		MetaID:       "130",
		MetaPlatform: "tmdb",
	}

	database.GetDB().Save(&person)

	credit := model.Credit{
		ItemID:     item.ID,
		PersonID:   person.ID,
		Department: "Acting",
		Role:       "Tony Stark",
	}

	database.GetDB().Create(&credit)

	// Display credits from item
	var foundItem model.Item

	database.GetDB().Preload("Credits").Preload("Credits.Person").First(&foundItem, item.ID)

	for _, credit := range foundItem.Credits {
		fmt.Printf("Movie: %s, Person: %s,Department: %s, Character: %s\n", foundItem.Title, credit.Person.Name, credit.Department, credit.Role)
	}
}
