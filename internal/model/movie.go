package model

import (
	db "ovo-server/internal/database"
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	TmdbID           uint      `json:"tmdb_id" gorm:"not null"`
	Title            string    `json:"title" gorm:"not null"`
	OriginalTitle    string    `json:"original_title" gorm:"not null"`
	Description      string    `json:"description" gorm:"not null"`
	ReleaseDate      time.Time `json:"release_date" gorm:"not null"`
	PosterPath       string    `json:"poster_path" gorm:"not null"`
	FilePath         string    `json:"file_path" gorm:"not null"`
	LastMetadataScan time.Time `json:"last_metadata_scan"`
}

func (movie *Movie) Save() error {

	// Save movie to database
	transaction := db.GetDB().Save(&movie)

	if transaction.Error != nil {
		return transaction.Error
	}

	return nil
}
