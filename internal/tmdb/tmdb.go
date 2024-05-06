package tmdb

import (
	"log"
	"ovo-server/internal/config"
	"strconv"
	"time"

	tmdbApi "github.com/ryanbradynd05/go-tmdb"
)

type TMDBMetadataItem struct {
	TmdbID         string
	Title          string
	OriginalTitle  string
	Description    string
	Tagline        string
	ReleaseDate    time.Time
	PosterPath     string
	BackdropPath   string
	SeasonNumber   int
	EpisodeNumber  int
	EpisodeTitle   string
	EpisodeAirDate string
}

type TMDBCredit struct {
	ItemTmdbID   string
	ItemType     string
	PersonTmdbID string
	Department   string
	Role         string
}

type TMDBPerson struct {
	TmdbID       string
	Name         string
	Biography    string
	Birthday     time.Time
	Deathday     time.Time
	PlaceOfBirth string
	ProfilePath  string
}

var api *tmdbApi.TMDb

func Init() {
	config := tmdbApi.Config{
		APIKey:   config.Variables.TMDBApiKey,
		Proxies:  nil,
		UseProxy: false,
	}

	api = tmdbApi.Init(config)

}

func GetMovieDetails(id int) *TMDBMetadataItem {
	details, err := api.GetMovieInfo(id, nil)
	if err != nil {
		log.Printf("Error getting movie details for id '%d': %s", id, err)
		return nil
	}

	metadata := &TMDBMetadataItem{
		TmdbID:        strconv.Itoa(details.ID),
		Title:         details.Title,
		OriginalTitle: details.OriginalTitle,
		Description:   details.Overview,
		PosterPath:    details.PosterPath,
		BackdropPath:  details.BackdropPath,
		Tagline:       details.Tagline,
	}

	releaseDate, err := time.Parse("2006-01-02", details.ReleaseDate)
	if err != nil {
		log.Printf("Error parsing release date for movie '%s': Received release date is: %s. \nError: %s. \nWon't get modified.", details.Title, details.ReleaseDate, err)
		return metadata
	}

	metadata.ReleaseDate = releaseDate
	return metadata
}

func SearchMovieByNameAndYear(name string, year int) *TMDBMetadataItem {
	options := make(map[string]string)
	if year != 0 {
		options["year"] = strconv.Itoa(year)
	}

	details, err := api.SearchMovie(name, options)
	if err != nil {
		log.Printf("Error searching movie '%s' with year '%d': %s", name, year, err)
		return nil
	}

	if len(details.Results) == 0 {
		log.Println("No movie found for ", name, " with year ", year)
		log.Println("Options: ", options)
		return nil
	}

	return GetMovieDetails(details.Results[0].ID)
}

// Finds TMDB ID from external meta provider. Using a platform and a ID, returns a ID.
// User is expected to use the returned ID on correct field (movie, show, season...)
// Types are
// - Movie
// - Show
// - Season
// - Episode
// - Person
func GetIDFromExternal(MetaProvider string, MetaID string) (int, error) {
	options := make(map[string]string)
	response, err := api.GetFind(MetaID, MetaProvider+"_id", options)
	var resultID int

	if err != nil {
		return 0, err
	}

	if len(response.MovieResults) > 0 {
		resultID = response.MovieResults[0].ID
	}
	if len(response.TvResults) > 0 {
		resultID = response.TvResults[0].ID
	}
	if len(response.TvSeasonResults) > 0 {
		resultID = response.TvSeasonResults[0].ID
	}
	if len(response.TvEpisodeResults) > 0 {
		resultID = response.TvEpisodeResults[0].ID
	}
	if len(response.PersonResults) > 0 {
		resultID = response.PersonResults[0].ID
	}

	return resultID, nil
}

func SearchMovie(name string) *TMDBMetadataItem {
	return SearchMovieByNameAndYear(name, 0)
}

func GetShowDetails(id int) *TMDBMetadataItem {
	details, err := api.GetTvInfo(id, nil)
	if err != nil {
		log.Printf("Error getting show details for id '%d': %s", id, err)
		return nil
	}

	metadata := &TMDBMetadataItem{
		TmdbID:        strconv.Itoa(details.ID),
		Title:         details.Name,
		OriginalTitle: details.OriginalName,
		Description:   details.Overview,
		PosterPath:    details.PosterPath,
		BackdropPath:  details.BackdropPath,
	}

	firstAiredDate, err := time.Parse("2006-01-02", details.FirstAirDate)

	if err != nil {
		log.Printf("Error parsing first aired date for show '%s': Received first aired date is: %s. \nError: %s. \nWon't get modified.", details.Name, details.FirstAirDate, err)
		return metadata
	}

	metadata.ReleaseDate = firstAiredDate
	return metadata
}

func SearchShowByNameAndYear(name string, year int) *TMDBMetadataItem {
	options := make(map[string]string)

	if year != 0 {
		options["year"] = strconv.Itoa(year)
	}

	details, err := api.SearchTv(name, options)
	if err != nil {
		log.Printf("Error searching show '%s': %s", name, err)
		return nil
	}

	if len(details.Results) == 0 {
		log.Println("No show found for ", name)
		return nil
	}

	return GetShowDetails(details.Results[0].ID)
}

func SearchShow(name string) *TMDBMetadataItem {
	return SearchShowByNameAndYear(name, 0)
}

func GetMovieCredits(id int) ([]TMDBCredit, error) {
	var credits []TMDBCredit
	var options = make(map[string]string)
	responseCredits, err := api.GetMovieCredits(id, options)
	if err != nil {
		log.Printf("Error getting movie credits for id '%d': %s", id, err)
		return nil, err
	}

	for _, credit := range responseCredits.Cast {
		credits = append(credits, TMDBCredit{
			ItemTmdbID:   strconv.Itoa(id),
			PersonTmdbID: strconv.Itoa(credit.ID),
			Department:   "cast",
			Role:         credit.Character,
		})
	}

	for _, credit := range responseCredits.Crew {
		credits = append(credits, TMDBCredit{
			ItemTmdbID:   strconv.Itoa(id),
			PersonTmdbID: strconv.Itoa(credit.ID),
			Department:   credit.Department,
			Role:         credit.Job,
		})
	}

	return credits, nil
}

func GetPerson(id string) (*TMDBPerson, error) {
	var options = make(map[string]string)
	personID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error converting id to int: %s", id)
		return nil, err
	}
	response, err := api.GetPersonInfo(personID, options)
	if err != nil {
		log.Printf("Error getting person info for id '%d': %s", personID, err)
		return nil, err
	}

	person := &TMDBPerson{
		TmdbID:       strconv.Itoa(response.ID),
		Name:         response.Name,
		Biography:    response.Biography,
		PlaceOfBirth: response.PlaceOfBirth,
		ProfilePath:  response.ProfilePath,
	}

	birthday, err := time.Parse("2006-01-02", response.Birthday)
	if err == nil {
		person.Birthday = birthday
	}

	deathday, err := time.Parse("2006-01-02", response.Deathday)
	if err == nil {
		person.Deathday = deathday
	}

	return person, nil
}
