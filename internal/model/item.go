package model

import (
	"log"
	db "ovo-server/internal/database"
	"ovo-server/internal/file"
	"ovo-server/internal/tmdb"
	"strconv"
	"time"

	"gorm.io/gorm"
)

const (
	ItemTypeMovie   = "Movie"
	ItemTypeShow    = "Show"
	ItemTypeSeason  = "Season"
	ItemTypeEpisode = "Episode"
)

type Item struct {
	gorm.Model
	LibraryID        uint      `json:"library" gorm:"not null"`
	ItemType         string    `gorm:"enum:show,season,episode,movie"`
	TmdbID           string    `json:"tmdb_id"`
	Title            string    `json:"title" gorm:"not null"`
	OriginalTitle    string    `json:"original_title" gorm:"not null"`
	Description      string    `json:"description"`
	ReleaseDate      time.Time `json:"release_date"`
	PosterPath       string    `json:"poster_path"`
	FilePath         string    `json:"file_path" gorm:"not null"`
	LastMetadataScan time.Time `json:"last_metadata_scan"`
	ParentItem       uint      `json:"parent_item"` // Show or Season ID
}

type ItemMetadata struct {
	TmdbID        string    `json:"tmdb_id"`
	Title         string    `json:"title"`
	OriginalTitle string    `json:"original_title"`
	Description   string    `json:"description"`
	ReleaseDate   time.Time `json:"release_date"`
	PosterPath    string    `json:"poster_path"`
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

func (item *Item) FetchMetadata() {
	var metadata *tmdb.TMDBMetadataItem

	switch item.ItemType {
	case ItemTypeMovie:
		if item.TmdbID != "" {
			tmdbID, _ := strconv.Atoi(item.TmdbID)
			metadata = tmdb.GetMovieDetails(tmdbID)
		} else if year := file.ParseYearFromFilename(item.FilePath); year != 0 {
			metadata = tmdb.SearchMovieByNameAndYear(item.Title, year)
		} else {
			metadata = tmdb.SearchMovie(item.Title)
		}

	case ItemTypeShow:
		if item.TmdbID != "" {
			tmdbID, _ := strconv.Atoi(item.TmdbID)
			metadata = tmdb.GetShowDetails(tmdbID)
		} else if year := file.ParseYearFromFilename(item.FilePath); year != 0 {
			metadata = tmdb.SearchShowByNameAndYear(item.Title, year)
		} else {
			metadata = tmdb.SearchShow(item.Title)
		}
	}

	if metadata != nil {
		log.Println("Updating metadata for", item.Title, "with ID", item.ID, "and tmdb ID", metadata.TmdbID)
		item.UpdateMovieMetadata(*metadata)
	}
}

func (item *Item) UpdateMovieMetadata(metadata tmdb.TMDBMetadataItem) {
	item.TmdbID = metadata.TmdbID
	item.Title = metadata.Title
	item.OriginalTitle = metadata.OriginalTitle
	item.Description = metadata.Description
	item.ReleaseDate = metadata.ReleaseDate
	item.PosterPath = metadata.PosterPath
	item.LastMetadataScan = time.Now()
	item.Save()
}
