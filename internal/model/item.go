package model

import (
	db "ovo-server/internal/database"
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	LibraryID        uint      `json:"library" gorm:"not null"`
	ItemType         string    `gorm:"enum:show,season,episode,movie"`
	TmdbID           uint      `json:"tmdb_id"`
	Title            string    `json:"title" gorm:"not null"`
	OriginalTitle    string    `json:"original_title" gorm:"not null"`
	Description      string    `json:"description"`
	ReleaseDate      time.Time `json:"release_date" gorm:"not null"`
	PosterPath       string    `json:"poster_path"`
	FilePath         string    `json:"file_path" gorm:"not null"`
	LastMetadataScan time.Time `json:"last_metadata_scan"`
	ParentItem       uint      `json:"parent_item"` // Show or Season ID
}

func (item *Item) Save() error {

	// Save movie to database
	transaction := db.GetDB().Save(&item)

	if transaction.Error != nil {
		return transaction.Error
	}

	return nil
}

func (item *Item) Delete() {
	db.GetDB().Delete(&item)
}
