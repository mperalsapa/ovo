package file

import (
	"errors"
	"regexp"
	"strconv"
)

func ParseSeasonDirname(dirname string) (int, error) {
	// We search for numbers in the dirname. This is expected to be the season number, e.g. "Season 1", "S22" or "50"
	seasonNumber := regexp.MustCompile(`\d+`).FindString(dirname)
	// If we couldn't find any numbers, we return an error
	if seasonNumber == "" {
		return 0, errors.New("could not parse season number")
	}

	// We convert the found number to an integer
	seasonNumberInt, err := strconv.Atoi(seasonNumber)
	// If we couldn't convert the number to an integer, we return an error
	// Although this should never happen, as we only search for numbers
	if err != nil {
		return 0, err
	}

	// We return the season number as an integer
	return seasonNumberInt, nil
}

func ParseEpisodeFilename(filename string) (int, error) {

	// first we remove the file extension
	filenameWithoutExtension := regexp.MustCompile(`(.+?)(\.[^.]*$|$)`).FindStringSubmatch(filename)[1]

	// We search for numbers in the filename. This would be like the season, but an episode also contains its season number.
	// So we expect the episode filename to be "Episode S1E1", that would be Episode 1 from Season 1.
	episodeNumber := regexp.MustCompile(`\d+`).FindAllString(filenameWithoutExtension, -1)
	// If we couldn't find any numbers, we return an error
	if episodeNumber == nil {
		return 0, errors.New("could not parse episode number")
	}

	// We convert the found numbers to integers
	episodeNumberInt, err := strconv.Atoi(episodeNumber[len(episodeNumber)-1])
	// If we couldn't convert the number to an integer, we return an error
	// Although this should never happen, as we only search for numbers
	if err != nil {
		return 0, err
	}

	// We return the episode number as an integer
	return episodeNumberInt, nil
}
