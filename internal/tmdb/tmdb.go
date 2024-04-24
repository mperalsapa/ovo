package tmdb

import (
	"log"
	"ovo-server/internal/config"
	"ovo-server/internal/file"
	"strconv"

	tmdbApi "github.com/ryanbradynd05/go-tmdb"
)

var api *tmdbApi.TMDb

func Init() {
	config := tmdbApi.Config{
		APIKey:   config.Variables.TMDBApiKey,
		Proxies:  nil,
		UseProxy: false,
	}

	api = tmdbApi.Init(config)

}

func FindMovieByFileInfoList(fileInfoList []file.FileMetaInfo) []*tmdbApi.Movie {
	movies := []*tmdbApi.Movie{}

	for _, fileInfo := range fileInfoList {
		movie := FindMovieByFileInfo(fileInfo)
		if movie != nil {
			movies = append(movies, movie)
		}
	}

	return movies

}

func FindMovieByFileInfo(fileInfo file.FileMetaInfo) *tmdbApi.Movie {
	if fileInfo.MetaProvider != "" && fileInfo.MetaId != "" {
		intId, err := strconv.Atoi(fileInfo.MetaId)
		if err != nil {
			log.Printf("Error converting MetaId to int on %s: %s", fileInfo.MetaId, err)
			return nil
		}
		return GetMovieDetails(intId)
	}

	if fileInfo.Name != "" && fileInfo.Year != 0 {
		result := SearchMovieByNameAndYear(fileInfo.Name, fileInfo.Year)

		if result != nil {
			return GetMovieDetails(result.ID)
		}

		log.Printf("No movie found for '%s' with year '%d'", fileInfo.Name, fileInfo.Year)
		return nil
	}

	if fileInfo.Name != "" {
		result := SearchMovie(fileInfo.Name)

		if result != nil {
			return GetMovieDetails(result.ID)
		}

		log.Printf("No movie found for '%s'", fileInfo.Name)
		return nil
	}
	return nil
}

func GetMovieDetails(id int) *tmdbApi.Movie {
	details, err := api.GetMovieInfo(id, nil)
	if err != nil {
		log.Printf("Error getting movie details for id '%d': %s", id, err)
		return nil
	}

	return details
}

func SearchMovieByNameAndYear(name string, year int) *tmdbApi.MovieShort {
	options := make(map[string]string)
	options["year"] = strconv.Itoa(year)
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

	return &details.Results[0]
}

func SearchMovie(name string) *tmdbApi.MovieShort {
	log.Println("Searching Movie by name only: ", name)
	var options map[string]string
	details, err := api.SearchMovie(name, options)
	if err != nil {
		log.Printf("Error searching movie '%s': %s", name, err)
		return nil
	}

	if len(details.Results) == 0 {
		return nil
	}

	return &details.Results[0]
}
