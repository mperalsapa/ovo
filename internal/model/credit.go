package model

import (
	"log"
	"ovo-server/internal/database"
	"ovo-server/internal/tmdb"
	"strconv"

	"gorm.io/gorm"
)

type Credit struct {
	gorm.Model
	ItemID     uint   `json:"item_id" gorm:"not null"`
	PersonID   uint   `json:"person_id" gorm:"not null"`
	Person     Person `json:"person"`
	Department string `json:"department"`
	Role       string `json:"role"`
}

func FetchCredits(item Item) {
	// var currentCredits []Credit
	// database.GetDB().Preload("Credits").Preload("Credits.Person").Find(&currentCredits, item.ID)
	var credits []tmdb.TMDBCredit
	var err error

	itemID, err := strconv.Atoi(item.MetaID)
	if err != nil {
		log.Println(err)
		return
	}

	switch item.ItemType {
	case ItemTypeMovie:
		credits, err = tmdb.GetMovieCredits(itemID)
	}

	if err != nil {
		log.Println(err)
		return
	}

	for _, credit := range credits {
		var person Person

		database.GetDB().First(&person, credit.PersonTmdbID)
		if person.ID == 0 {
			personMeta, _ := tmdb.GetPerson(credit.PersonTmdbID)
			person = Person{
				Name:         personMeta.Name,
				MetaID:       credit.PersonTmdbID,
				MetaPlatform: "tmdb",
				Biography:    personMeta.Biography,
			}
		}

		// newCredit := Credit{
		// 	ItemID:    item.ID,
		// 	PersonID:  person.ID,
		// 	Character: credit.Character,
		// }
	}
}

func (credit *Credit) Save() error {
	transaction := database.GetDB().Save(&credit)

	if transaction.Error != nil {
		return transaction.Error
	}
	return nil
}
