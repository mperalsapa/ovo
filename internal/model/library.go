package model

import "gorm.io/gorm"

type LibraryType string

const (
	LibraryTypeMovie LibraryType = "movie"
	LibraryTypeShow  LibraryType = "show"
)

type Library struct {
	gorm.Model
	Type LibraryType `json:"type" gorm:"not null; enum('movie', 'show')"`
}
