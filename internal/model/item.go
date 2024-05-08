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
	Title            string    `json:"title" gorm:"not null;index"`
	OriginalTitle    string    `json:"original_title" gorm:"not null;index"`
	Description      string    `json:"description"`
	TagLine          string    `json:"tag_line" gorm:"default:null"`
	ReleaseDate      time.Time `json:"release_date"`
	PosterPath       string    `json:"poster_path" gorm:"default:null"`
	FilePath         string    `json:"file_path" gorm:"not null"`
	LastMetadataScan time.Time `json:"last_metadata_scan" gorm:"default:null"`
	ParentItem       uint      `json:"parent_item"` // Show or Season ID
	Credits          []Credit  `json:"credits" gorm:"constraint:OnDelete:CASCADE"`
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

func GetItemById(id uint) (Item, error) {
	var item Item
	result := db.GetDB().First(&item, id)

	if result.Error != nil {
		return item, result.Error
	}

	return item, nil
}

func (item *Item) FetchMetadata() {
	var metadata *tmdb.TMDBMetadataItem
	log.Println("Fetching metadata for", item.Title, "with ID", item.ID, "and tmdb ID", item.MetaID)
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
		} else if year := file.ParseYearFromFilename(item.FilePath); year != 0 {
			metadata = tmdb.SearchMovieByNameAndYear(item.Title, year)
		} else {
			metadata = tmdb.SearchMovie(item.Title)
		}

	case ItemTypeShow:
		if item.MetaID != "" {
			tmdbID, err := strconv.Atoi(item.MetaID)
			if err != nil {
				log.Printf("Error converting ID to int: %s", item.MetaID)
				return
			}
			metadata = tmdb.GetShowDetails(tmdbID)
		} else if year := file.ParseYearFromFilename(item.FilePath); year != 0 {
			metadata = tmdb.SearchShowByNameAndYear(item.Title, year)
		} else {
			metadata = tmdb.SearchShow(item.Title)
		}
	case ItemTypeSeason:
		var parentItem Item
		db.GetDB().First(&parentItem, item.ParentItem)
		if parentItem.MetaID == "" {
			log.Println("Parent item has no metadata. Skipping season metadata fetch.")
			return
		}

		tmdbID, err := strconv.Atoi(parentItem.MetaID)
		if err != nil {
			log.Printf("Error converting ID to int: %s", parentItem.MetaID)
			return
		}

		seasonNumber, err := strconv.Atoi(item.Title)
		if err != nil {
			log.Printf("Error converting season number to int: %s", item.Title)
			return
		}

		metadata = tmdb.GetSeasonDetails(tmdbID, seasonNumber)
	}

	if metadata == nil {
		return
	}

	log.Println("Updating metadata for", item.Title, "with ID", item.ID, "and tmdb ID", metadata.TmdbID)
	item.UpdateMovieMetadata(*metadata)

	item.FetchCredits()
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
	item.TagLine = metadata.Tagline
	item.Save()
}

func (item *Item) FetchCredits() {
	var credits []tmdb.TMDBCredit
	var err error

	// If item does not have meta id we can't know what to search for
	if item.MetaID == "" {
		return
	}

	itemID, err := strconv.Atoi(item.MetaID)
	if err != nil {
		log.Println(err)
		return
	}

	switch item.ItemType {
	case ItemTypeMovie:
		{
			credits, err = tmdb.GetMovieCredits(itemID)
		}
	case ItemTypeShow:
		{
			credits, err = tmdb.GetShowCredits(itemID)
		}
	}

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Processing credits for", item.Title, "with ID", item.ID, "and tmdb ID", item.MetaID)
	maxCastPerson := 18
	castCount := 0
	for _, credit := range credits {

		if credit.Department == "cast" {
			if castCount >= maxCastPerson {
				continue
			}
			castCount++
		}

		// Check if person already exists in database
		var person Person
		db.GetDB().Where(&Person{MetaID: credit.PersonTmdbID}).First(&person)
		if person.ID == 0 {
			personMeta, _ := tmdb.GetPerson(credit.PersonTmdbID)
			person = Person{
				Name:         personMeta.Name,
				MetaID:       credit.PersonTmdbID,
				MetaPlatform: "tmdb",
				ProfilePath:  personMeta.ProfilePath,
				Biography:    personMeta.Biography,
				PlaceOfBirth: personMeta.PlaceOfBirth,
				Birthday:     personMeta.Birthday,
				Deathday:     personMeta.Deathday,
			}
			err := person.Save()
			if err != nil {
				log.Printf("Error saving person %s: %s", person.Name, err)
				log.Println("Birthday:", person.Birthday)

				continue
			}
		}

		// Check if credit already exists in database
		var existingCredit Credit
		db.GetDB().Where(&Credit{ItemID: item.ID, PersonID: person.ID}).First(&existingCredit)
		if existingCredit.ID == 0 {
			newCredit := Credit{
				ItemID:     item.ID,
				PersonID:   person.ID,
				Department: credit.Department,
				Role:       credit.Role,
			}
			err := newCredit.Save()
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

}

func (item *Item) GetDirectors() []Credit {
	var credits []Credit
	db.GetDB().Where("item_id = ? AND department = ?", item.ID, "Directing").Preload("Person").Find(&credits)
	return credits
}

func (item *Item) GetWriters() []Credit {
	var credits []Credit
	db.GetDB().Where("item_id = ? AND department = ?", item.ID, "Writing").Preload("Person").Find(&credits)
	return credits
}

func (item *Item) GetCast() []Credit {
	var credits []Credit
	db.GetDB().Where("item_id = ? AND department = ?", item.ID, "cast").Preload("Person").Find(&credits)
	return credits
}
