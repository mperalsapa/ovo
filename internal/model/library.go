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
	Type  LibraryType `json:"type" form:"type" gorm:"not null; enum('movie', 'show')"`
	Name  string      `json:"name" form:"name" gorm:"not null"`
	Paths []string    `json:"paths" form:"paths[]" gorm:"serializer:json"`
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

	library.removeEmptyPaths()
	if len(library.Paths) == 0 {
		return errors.New("paths is required")
	}

	return nil
}

func (library *Library) removeEmptyPaths() {
	var paths []string
	for _, path := range library.Paths {
		if path != "" {
			paths = append(paths, path)
		}
	}
	library.Paths = paths
}

func GetLibraries() []Library {
	var libraries []Library
	db.GetDB().Find(&libraries)
	return libraries
}

func GetLibraryById(id uint) (Library, error) {
	var library Library
	transaction := db.GetDB().First(&library, id)
	if transaction.Error != nil {
		return Library{}, transaction.Error
	}
	return library, nil
}

func DeleteLibrary(id uint) error {
	library, err := GetLibraryById(id)
	if err != nil {
		return err
	}

	transaction := db.GetDB().Delete(&library)
	if transaction.Error != nil {
		return transaction.Error
	}
	return nil
}

func (library *Library) DeleteLibrary() {
	db.GetDB().Delete(&library)
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
