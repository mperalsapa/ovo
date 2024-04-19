package model

import (
	db "ovo-server/internal/database"

	"gorm.io/gorm"
)

type LibraryType string

const (
	LibraryTypeMovie LibraryType = "movie"
	LibraryTypeShow  LibraryType = "show"
)

type Library struct {
	gorm.Model
	Type  LibraryType `json:"type" gorm:"not null; enum('movie', 'show')"`
	Name  string      `json:"name" gorm:"not null"`
	Paths []string    `json:"paths" gorm:"serializer:json"`
}

func GetLibraries() []Library {
	var libraries []Library
	db.GetDB().Find(&libraries)
	return libraries
}
