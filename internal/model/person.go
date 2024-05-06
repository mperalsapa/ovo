package model

import (
	"ovo-server/internal/database"
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	MetaID       string    `json:"meta_id" gorm:"not null"`
	MetaPlatform string    `json:"meta_platform" gorm:"not null"`
	Name         string    `json:"name" gorm:"not null;index"`
	Biography    string    `json:"biography" gorm:"default:null"`
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
