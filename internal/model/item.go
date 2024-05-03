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
	ItemTypeMovie    = "Movie"
	ItemTypeShow     = "Show"
	ItemTypeSeason   = "Season"
	ItemTypeEpisode  = "Episode"
	MetaProviderTMDB = "tmdb"
)

type Item struct {
	gorm.Model
	LibraryID        uint      `json:"library" gorm:"not null"`
	ItemType         string    `gorm:"enum:show,season,episode,movie"`
	MetaProvider     string    `json:"meta_platform"`
	MetaID           string    `json:"meta_id"`
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
		if item.MetaID != "" {
			var tmdbID int
			var err error
			// If metadata is from external ID, we search for TMDB ID using external
			if item.MetaProvider != MetaProviderTMDB {
				log.Printf("Using external ID. Platform: %s, ID: %s", item.MetaProvider, item.MetaID)
				tmdbID, err = tmdb.GetIDFromExternal(item.MetaProvider, item.MetaID)
				if err != nil {
					log.Printf("Error searching by external ID: %s", err)
					return
				}
				log.Printf("ID result from %s in %s: %d", item.MetaID, item.MetaProvider, tmdbID)
			} else {
				tmdbID, err = strconv.Atoi(item.MetaID)
				if err != nil {
					log.Printf("Error converting ID to int: %s", item.MetaID)
					return
				}
			}

			metadata = tmdb.GetMovieDetails(tmdbID)
		}

		if metadata == nil {
			if year := file.ParseYearFromFilename(item.FilePath); year != 0 {
				metadata = tmdb.SearchMovieByNameAndYear(item.Title, year)
			} else {
				metadata = tmdb.SearchMovie(item.Title)
			}
		}

	case ItemTypeShow:
		if item.MetaID != "" {
			tmdbID, _ := strconv.Atoi(item.MetaID)
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
	item.MetaProvider = MetaProviderTMDB
	item.MetaID = metadata.TmdbID
	item.Title = metadata.Title
	item.OriginalTitle = metadata.OriginalTitle
	item.Description = metadata.Description
	item.ReleaseDate = metadata.ReleaseDate
	item.PosterPath = metadata.PosterPath
	item.LastMetadataScan = time.Now()
	item.Save()
}
