package model

import (
	"log"
	"ovo-server/internal/database"
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	MetaID       string    `json:"meta_id" gorm:"not null"`
	MetaPlatform string    `json:"meta_platform" gorm:"not null"`
	Name         string    `json:"name" gorm:"not null;index"`
	Biography    string    `json:"biography" gorm:"default:null; type:longtext"`
	Birthday     time.Time `json:"birthday" gorm:"default:null"`
	Deathday     time.Time `json:"deathday" gorm:"default:null"`
	PlaceOfBirth string    `json:"place_of_birth" gorm:"default:null"`
	ProfilePath  string    `json:"profile_path" gorm:"default:null"`
	Credits      []Credit  `json:"credits"`
}

func (person *Person) Save() error {
	transaction := database.GetDB().Save(&person)

	if transaction.Error != nil {
		return transaction.Error
	}

	return nil
}

func (person *Person) GetCredits() []Credit {
	var credits []Credit
	transaction := database.GetDB().Preload("Person").Where("person_id = ?", person.ID).Find(&credits)

	if transaction.Error != nil {
		log.Println("Error when getting credits for person: ", person.ID, transaction.Error)
		return nil
	}

	return credits
}

func (person *Person) LoadCredits() {
	person.Credits = person.GetCredits()
}

func GetPersonById(id uint) (*Person, error) {
	var person Person
	transaction := database.GetDB().First(&person, id)

	if transaction.Error != nil {
		return nil, transaction.Error
	}

	return &person, nil
}

func (p *Person) GetCreditItems() []Item {
	var credits []Credit
	var items []Item
	uniqueItems := make(map[uint]Item)

	transaction := database.GetDB().Distinct("item_id").Preload("Item").Where("person_id = ?", p.ID).Find(&credits)

	if transaction.Error != nil {
		log.Println("Error when getting person with credits and items: ", p.ID, transaction.Error)
		return nil
	}

	for _, c := range credits {
		uniqueItems[c.ItemID] = c.Item
	}

	log.Println(len(uniqueItems))

	for _, i := range uniqueItems {
		items = append(items, i)
	}

	log.Println(len(items))

	return items
}
