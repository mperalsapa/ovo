package model

import (
	"errors"
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

func (library *Library) Equals(other Library) bool {
	return library.ID == other.ID &&
		library.Type == other.Type &&
		library.Name == other.Name
}

func (library *Library) Validate() error {
	if library.Type == "" {
		return errors.New("type is required")
	}
	if library.Name == "" {
		return errors.New("name is required")
	}
	if len(library.Paths) == 0 {
		return errors.New("paths is required")
	}
	return nil
}

func GetLibraries() []Library {
	var libraries []Library
	db.GetDB().Find(&libraries)
	return libraries
}

func GetLibraryById(id uint) Library {
	var library Library
	db.GetDB().First(&library, id)
	return library
}

func (library *Library) SaveLibrary() error {
	if err := library.Validate(); err != nil {
		return err
	}
	transaction := db.GetDB().Save(&library)

	if transaction.Error != nil {
		return transaction.Error
	}

	return nil
}
