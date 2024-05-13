package model

import (
	"fmt"
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
	LibraryID        uint          `json:"library" gorm:"not null"`
	ItemType         string        `gorm:"enum:show,season,episode,movie"`
	MetaProvider     string        `json:"meta_platform"`
	MetaID           string        `json:"meta_id"`
	MetaRating       float32       `json:"meta_rating"`
	Title            string        `json:"title" gorm:"not null;index"`
	OriginalTitle    string        `json:"original_title" gorm:"not null;index"`
	Description      string        `json:"description"`
	TagLine          string        `json:"tag_line" gorm:"default:null"`
	ReleaseDate      time.Time     `json:"release_date"`
	EndDate          *time.Time    `json:"end_date"`
	PosterPath       string        `json:"poster_path" gorm:"default:null"`
	FilePath         string        `json:"file_path" gorm:"not null"`
	Duration         time.Duration `json:"duration"`
	LastMetadataScan *time.Time    `json:"last_metadata_scan"`
	ParentItem       uint          `json:"parent_item"` // Show or Season ID
	Credits          []Credit      `json:"credits" gorm:"constraint:OnDelete:CASCADE"`
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
	// log.Println("Fetching metadata for", item.Title, "with ID", item.ID, "and tmdb ID", item.MetaID)
	switch item.ItemType {
	case ItemTypeMovie:
		if item.MetaID != "" {
			var metaID int
			var err error
			// If metadata is from external ID, we search for TMDB ID using external
			if item.MetaProvider != MetaProviderTMDB {
				log.Printf("Using external ID. Platform: %s, ID: %s", item.MetaProvider, item.MetaID)
				metaID, err = tmdb.GetIDFromExternal(item.MetaProvider, item.MetaID)
				if err != nil {
					log.Printf("Error searching by external ID: %s", err)
					return
				}
			} else {
				metaID, err = strconv.Atoi(item.MetaID)
				if err != nil {
					log.Printf("Error converting ID to int: %s", item.MetaID)
					return
				}
			}

			metadata = tmdb.GetMovieDetails(metaID)
		} else if year := file.ParseYearFromFilename(item.FilePath); year != 0 {
			metadata = tmdb.SearchMovieByNameAndYear(item.Title, year)
		} else {
			metadata = tmdb.SearchMovie(item.Title)
		}

	case ItemTypeShow:
		if item.MetaID != "" {
			metaID, err := strconv.Atoi(item.MetaID)
			if err != nil {
				log.Printf("Error converting ID to int: %s", item.MetaID)
				return
			}
			metadata = tmdb.GetShowDetails(metaID)
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

		showMetaID, err := strconv.Atoi(parentItem.MetaID)
		if err != nil {
			log.Printf("Error converting ID to int: %s", parentItem.MetaID)
			return
		}

		seasonNumber, err := strconv.Atoi(item.MetaID)
		if err != nil {
			log.Printf("Error converting season number to int: %s", item.Title)
			return
		}

		metadata = tmdb.GetSeasonDetails(showMetaID, seasonNumber)
	case ItemTypeEpisode:
		var seasonItem Item
		var showItem Item
		db.GetDB().First(&seasonItem, item.ParentItem)
		db.GetDB().First(&showItem, seasonItem.ParentItem)

		if seasonItem.MetaID == "" || showItem.MetaID == "" {
			log.Println("Parent item has no metadata. Skipping episode metadata fetch.")
			return
		}

		showMetaID, err := strconv.Atoi(showItem.MetaID)
		if err != nil {
			log.Printf("Error converting ID to int: %s", showItem.MetaID)
			return
		}

		seasonNumber, err := strconv.Atoi(seasonItem.MetaID)
		if err != nil {
			log.Printf("Error converting season number to int: %s", seasonItem.MetaID)
			return
		}

		episodeNumber, err := strconv.Atoi(item.MetaID)
		if err != nil {
			log.Printf("Error converting episode number to int: %s", item.MetaID)
			return
		}

		metadata = tmdb.GetEpisodeDetails(showMetaID, seasonNumber, episodeNumber)

	}

	if metadata == nil {
		return
	}

	log.Println("Updating metadata for", item.Title, "with ID", item.ID, "and tmdb ID", metadata.TmdbID)

	// log.Println("Updating metadata for", item.Title, "with ID", item.ID, "and tmdb ID", metadata.TmdbID)
	item.UpdateMovieMetadata(*metadata)

	item.FetchCredits()
}

func (item *Item) UpdateMovieMetadata(metadata tmdb.TMDBMetadataItem) {
	// Metdata from external source
	item.MetaProvider = MetaProviderTMDB
	item.MetaID = metadata.TmdbID
	item.MetaRating = metadata.Rating
	item.Title = metadata.Title
	item.OriginalTitle = metadata.OriginalTitle
	item.Description = metadata.Description
	item.ReleaseDate = metadata.ReleaseDate
	item.PosterPath = metadata.PosterPath
	item.TagLine = metadata.Tagline

	if metadata.EndDate != (time.Time{}) {
		item.EndDate = &metadata.EndDate
	}

	now := time.Now()
	item.LastMetadataScan = &now

	err := item.Save()
	if err != nil {
		log.Println("Error saving item", item.Title, "with ID", item.ID, "and tmdb ID", item.MetaID, ":", err)
	}
}

func (item *Item) UpdateRuntime() {
	itemFileMetadata := file.GetFileMetadata(item.FilePath)
	item.Duration = itemFileMetadata.Duration()
	err := item.Save()
	if err != nil {
		log.Println("Error saving item", item.Title, "with ID", item.ID, "and tmdb ID", item.MetaID, ":", err)
	}
}

func (item *Item) UpdateItemRuntime() {
	if item.ItemType == ItemTypeMovie || item.ItemType == ItemTypeEpisode {
		item.UpdateRuntime()
		return
	}

	var children []Item
	db.GetDB().Where("parent_item = ?", item.ID).Find(&children)
	for _, child := range children {
		child.UpdateItemRuntime()
	}

	var totalDuration time.Duration
	db.GetDB().Model(&children).Select("SUM(duration)").Where(&Item{ParentItem: item.ID}).Row().Scan(&totalDuration)
	item.Duration = totalDuration

	err := item.Save()
	if err != nil {
		log.Println("Error saving item", item.Title, "with ID", item.ID, "and tmdb ID", item.MetaID, ":", err)
	}

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

		// Getting only credits within role of the next list
		// - Director
		// - Writer
		// - cast
		if !(credit.Department == "Directing" && credit.Role == "Director") &&
			!(credit.Department == "Writing") &&
			!(credit.Department == "cast") {
			continue
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
	db.GetDB().Where("item_id = ? AND department = ? AND role = ?", item.ID, "Directing", "Director").Preload("Person").Find(&credits)
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

func (item *Item) GetChildren(itemType string) []Item {
	var children []Item
	db.GetDB().Where("parent_item = ? AND item_type = ?", item.ID, itemType).Find(&children)
	return children
}

func (item *Item) GetSeasons() []Item {
	return item.GetChildren(ItemTypeSeason)
}

func (item *Item) GetEpisodes() []Item {
	return item.GetChildren(ItemTypeEpisode)
}

func (item *Item) GetFancyDuration() string {
	var durationString string

	hours := int(item.Duration.Hours())
	if hours > 0 {
		durationString = fmt.Sprintf("%dh ", hours)
	}

	minutes := int(item.Duration.Minutes()) - hours*60
	if minutes > 0 {
		durationString = durationString + fmt.Sprintf("%dm", minutes)
	}

	return durationString
}
