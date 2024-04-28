package model

import (
	db "ovo-server/internal/database"
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	LibraryID        uint      `json:"library_id"` // Library ID
	ItemType         string    `gorm:"enum:show,season,episode,movie"`
	TmdbID           uint      `json:"tmdb_id" gorm:"not null"`
	Title            string    `json:"title" gorm:"not null"`
	OriginalTitle    string    `json:"original_title" gorm:"not null"`
	Description      string    `json:"description" gorm:"not null"`
	ReleaseDate      time.Time `json:"release_date" gorm:"not null"`
	PosterPath       string    `json:"poster_path" gorm:"not null"`
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
