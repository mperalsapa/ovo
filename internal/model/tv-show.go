package model

import (
	"time"

	"gorm.io/gorm"
)

type Show struct {
	gorm.Model
	TmdbID        uint      `json:"tmdb_id" gorm:"not null"`
	Title         string    `json:"title" gorm:"not null"`
	OriginalTitle string    `json:"original_title" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null"`
	ReleaseDate   time.Time `json:"release_date" gorm:"not null"`
	PosterPath    string    `json:"poster_path" gorm:"not null"`
	Seasons       []Season
}

type Season struct {
	gorm.Model
	SeasonNumber uint `json:"season_number" gorm:"not null"`
	Episodes     []Episode
}

type Episode struct {
	gorm.Model
	TmdbID uint `json:"tmdb_id" gorm:"not null"`
}
