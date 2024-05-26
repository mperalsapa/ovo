package model

import (
	"errors"
	"log"
	"math/rand"
	db "ovo-server/internal/database"
	"ovo-server/internal/file"
	"ovo-server/internal/tmdb"
	"path"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LibraryType string

const (
	LibraryTypeMovie = "movie"
	LibraryTypeShow  = "show"
)

type Library struct {
	gorm.Model
	Type      string   `json:"type" form:"type" gorm:"not null; enum('movie', 'show')"`
	Name      string   `json:"name" form:"name" gorm:"not null"`
	ImagePath string   `json:"image_path" form:"image_path"`
	Paths     []string `json:"paths" form:"paths[]" gorm:"serializer:json"`
	Items     []Item   `json:"items"`
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
	for _, libPath := range library.Paths {
		if libPath != "" {
			paths = append(paths, libPath)
		}
	}
	library.Paths = paths
}

// This function removes duplicated Items based on his path
func (library *Library) DeDuplicateItems() {
	UniqueFilePaths := make(map[string]bool)
	for _, item := range library.Items {
		// Verificar si el FilePath del item ya está en UniqueFilePaths
		if UniqueFilePaths[item.FilePath] {
			log.Println("Deleting duplicated item: ", item.Title)
			item.Delete()

		} else {
			// Si no está en UniqueFilePaths, agregar el FilePath a UniqueFilePaths
			UniqueFilePaths[item.FilePath] = true
		}
	}

	library.Items = nil
	library.GetItems("")
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

func (library *Library) Save() error {
	if err := library.Validate(); err != nil {
		return err
	}
	transaction := db.GetDB().Save(&library)

	if transaction.Error != nil {
		return transaction.Error
	}

	return nil
}

func SanitizeOrderBy(order string) string {
	var newOrder string
	fields := strings.Split(order, " ")

	if len(fields) == 0 {
		return ""
	}

	switch fields[0] {
	case "title":
		newOrder = "title"
	case "release_date":
		newOrder = "release_date"
	case "created_at":
		newOrder = "created_at"
	case "duration":
		newOrder = "duration"
	case "meta_rating":
		newOrder = "meta_rating"
	default:
		newOrder = ""
	}

	if len(fields) < 2 {
		return newOrder
	}

	switch fields[1] {
	case "asc":
		newOrder += " asc"
	case "desc":
		newOrder += " desc"
	}

	return newOrder
}

// GetItems return all items from library
func (library *Library) GetItems(order string) []Item {
	var items []Item

	order = SanitizeOrderBy(order)
	if order == "" {
		db.GetDB().Where(Item{LibraryID: library.ID}).Find(&items)
	} else {
		db.GetDB().Order(order).Where(Item{LibraryID: library.ID}).Find(&items)
	}
	return items
}

// Loads all items from the library into the Items field
func (library *Library) LoadItems(orderBy string) {
	library.Items = library.GetItems(orderBy)
}

func (library *Library) ScanLibrary() error {
	// checking if library has paths
	if len(library.Paths) == 0 {
		return errors.New("no paths to scan")
	}

	// Getting current items from database
	library.LoadItems("")

	// Scanning for new items that are not in the database yet
	library.ScanForNewItems()

	// Deduplicate items
	library.DeDuplicateItems()

	// Remove orphan items
	library.RemoveOrphanItems()

	// Fetch metadata for all items
	for _, item := range library.GetItems("") {
		item.FetchMetadata()
	}

	// Separately fetch runtime for files. This is done separately because when updating runtime
	// we make queries to the database, and that data is not available back into the
	// library items list without reloading them.
	for _, item := range library.GetItems("") {
		// Fetch runtime from file
		if item.ItemType == ItemTypeMovie || item.ItemType == ItemTypeShow {
			item.UpdateItemRuntime()
		}
	}

	// Finally, we get a random item from the library to set it as the library image
	// We do 10 tries to get a random item with an image, if we don't get one, we don't set the image
	for i := 0; i < min(10, len(library.Items)); i++ {
		imageItem := library.Items[rand.Intn(len(library.Items))]
		itemID, err := strconv.Atoi(imageItem.MetaID)
		if err != nil {
			continue
		}

		if imageItem.ItemType == ItemTypeMovie {
			imagePath := tmdb.GetMovieBackdrop(uint(itemID))
			library.ImagePath = imagePath
		} else if imageItem.ItemType == ItemTypeShow {
			imagePath := tmdb.GetShowBackdrop(uint(itemID))
			library.ImagePath = imagePath
		}

		// If we got an image, we break the loop
		if library.ImagePath != "" {
			break
		}

	}

	library.Save()

	return nil
}

// ItemExistsOnDisk checks if the item exists on disk
func (library *Library) ItemExistsOnDisk(item Item) bool {
	for range library.Paths {
		if item.FilePath != "" {
			if file.Exists(item.FilePath) {
				return true
			}
		}
	}
	return false
}

func (library *Library) GetItemByPath(path string) (item Item) {
	db.GetDB().Where(&Item{FilePath: path, LibraryID: library.ID}).First(&item)
	return item
}

// RemoveOrphanItems removes items that are not on disk anymore
func (library *Library) RemoveOrphanItems() {
	library.LoadItems("")
	for _, item := range library.Items {
		if !library.ItemExistsOnDisk(item) {
			item.Delete()
		}
	}
}

// ScanForNewItems scans the library paths for new items and adds them to the database as items without metadata
// Depending on the library type, it will scan for movies or shows
func (library *Library) ScanForNewItems() {
	switch library.Type {
	case LibraryTypeMovie:
		library.ScanForNewMovies()
	case LibraryTypeShow:
		library.ScanForNewShows()
	}
}

func (library *Library) GetPathsMap() map[string]bool {
	paths := make(map[string]bool)
	for _, item := range library.Items {
		paths[item.FilePath] = true
	}
	return paths
}

// Scans the given library paths for new movies.
// In case the movie is not in the database yet, it will be added
// The current items paths are what determines if that movie is already in the database.
// This will make that when a movie is renamed, it will be re-added as a new one, because it's path changed.
// Maybe in a future we could add a hash to the items to check if the file is the same, although this could be a bit expensive.
// TODO: Improve scanning funcionality using hashes or other methods to check if the file is the same
func (library *Library) ScanForNewMovies() {
	// Storing current paths in a map for faster lookup to check if files are already in the database
	currentPaths := library.GetPathsMap()

	for _, libPath := range library.Paths {
		files := file.ScanFiles(libPath)
		for _, movie := range files {
			log.Println("Movie detected:", movie)
			filePath := path.Join(libPath, movie)
			if currentPaths[filePath] {
				log.Println("Movie already in database:", movie)
				continue
			}
			fileInfo := file.ParseFilename(movie)
			item := Item{
				LibraryID:     library.ID,
				MetaProvider:  fileInfo.MetaProvider,
				MetaID:        fileInfo.MetaID,
				ItemType:      ItemTypeMovie,
				Title:         fileInfo.Name,
				OriginalTitle: fileInfo.Name,
				ReleaseDate:   time.Now(),
				FilePath:      filePath,
			}
			log.Println("Adding movie to database:", movie)
			err := item.Save()
			if err != nil {
				log.Println("Error saving movie:", err)
				continue
			}
			currentPaths[filePath] = true
		}
	}
}

func (library *Library) ScanForNewShows() {
	// Storing current paths in a map for faster lookup to check if files are already in the database
	currentPaths := library.GetPathsMap()
	for _, libPath := range library.Paths {
		shows := file.ScanDirectories(libPath)
		for _, show := range shows {
			log.Println("Show detected:", show)
			showPath := path.Join(libPath, show)
			var showItem Item
			// We check if current show exists in database. In case it doesn't, we create a dummy item containing basic info
			// about the show, and we store the item for later usage in seasons.
			if !currentPaths[showPath] {
				parsedShow := file.ParseFilename(show)
				showItem = Item{
					LibraryID:     library.ID,
					MetaProvider:  parsedShow.MetaProvider,
					MetaID:        parsedShow.MetaID,
					ItemType:      ItemTypeShow,
					Title:         parsedShow.Name,
					OriginalTitle: parsedShow.Name,
					ReleaseDate:   time.Now(),
					FilePath:      showPath,
				}
				showItem.Save()
			}

			// Because we dont know if the item is new or existing, we need to retrieve it from the database.
			// This will be used to reference seasons to their show.
			showItem = library.GetItemByPath(showPath)
			library.ScanForNewSeasons(showItem, currentPaths)
		}
	}
}

func (library *Library) ScanForNewSeasons(show Item, itemsPaths map[string]bool) {

	if itemsPaths == nil {
		itemsPaths = library.GetPathsMap()
	}

	// Start season scan
	seasons := file.ScanDirectories(show.FilePath)
	for _, season := range seasons {
		log.Println("Season detected:", season)
		seasonPath := path.Join(show.FilePath, season)
		var seasonItem Item
		if !itemsPaths[seasonPath] {
			parsedSeason, err := file.ParseSeasonDirname(season)
			if err != nil {
				log.Printf("Error parsing season %s: %s", seasonPath, err)
				continue
			}
			seasonItem = Item{
				LibraryID:   library.ID,
				ItemType:    ItemTypeSeason,
				MetaID:      strconv.Itoa(parsedSeason),
				Title:       strconv.Itoa(parsedSeason),
				ReleaseDate: time.Now(),
				FilePath:    seasonPath,
				ParentItem:  show.ID,
			}
			seasonItem.Save()
		}
		seasonItem = library.GetItemByPath(seasonPath)
		library.ScanForNewEpisodes(seasonItem, itemsPaths)
	}
}

func (library *Library) ScanForNewEpisodes(season Item, itemsPaths map[string]bool) {
	if itemsPaths == nil {
		itemsPaths = library.GetPathsMap()
	}

	episodes := file.ScanFiles(season.FilePath)
	for _, episode := range episodes {
		log.Println("Episode detected:", episode)
		episodePath := path.Join(season.FilePath, episode)
		var episodeItem Item
		if !itemsPaths[episodePath] {
			parsedEpisode, err := file.ParseEpisodeFilename(episode)
			if err != nil {
				log.Printf("Error parsing episode %s: %s", episodePath, err)
				continue
			}
			episodeItem = Item{
				LibraryID:   library.ID,
				ItemType:    ItemTypeEpisode,
				MetaID:      strconv.Itoa(parsedEpisode),
				Title:       strconv.Itoa(parsedEpisode),
				ReleaseDate: time.Now(),
				FilePath:    episodePath,
				ParentItem:  season.ID,
			}
			episodeItem.Save()
		}
	}
}

func (library *Library) GetLibraryMainItems() []Item {
	mainItems := []Item{}

	for _, item := range library.Items {
		if item.ItemType == ItemTypeMovie || item.ItemType == ItemTypeShow {
			mainItems = append(mainItems, item)
		}
	}

	return mainItems
}
