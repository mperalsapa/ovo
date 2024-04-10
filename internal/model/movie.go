package model

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	TmdbID        uint      `json:"tmdb_id" gorm:"not null"`
	Title         string    `json:"title" gorm:"not null"`
	OriginalTitle string    `json:"original_title" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null"`
	ReleaseDate   time.Time `json:"release_date" gorm:"not null"`
	PosterPath    string    `json:"poster_path" gorm:"not null"`
	FilePath      string    `json:"file_path" gorm:"not null"`
}
