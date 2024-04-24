package file

import (
	"log"
	"os"
	"regexp"
	"strconv"
)

type FileInfo struct {
	Name string
	Path string
	Size int64
}

type FileMetaInfo struct {
	Name         string
	Year         int
	MetaProvider string
	MetaId       string
}

// This function scans a dir and returns a slice of files
func ScanPath(path string) []string {
	// entries, err := ioutil.ReadDir(path)
	entries, err := os.ReadDir(path)
	files := []string{}
	if err != nil {
		log.Println(err)
	}

	for _, file := range entries {
		if file.IsDir() {
			continue
		}

		files = append(files, file.Name())
	}
	return files
}

// This function tries to get metadata from the given filename.
// It expects names containing name, and optionally year in parentheses and meta provider inside brackets followed by its id.
// followed by its provider (e.g. [tmdb-1234])
// The name is expected to be separated by dots, spaces or underscores
// Example 1: "The Matrix (1999)"
// Example 2: "The_Matrix_(1999)_[tmdb-603]"
// Example 3: "The Matrix 1999"
func ParseFilename(filename string) FileMetaInfo {
	fileInfo := FileMetaInfo{}
	// removing file extension
	filenameWithoutExtension := regexp.MustCompile(`(.+?)(\.[^.]*$|$)`).FindStringSubmatch(filename)[1]

	metaId := regexp.MustCompile(`\[(.+)-(.+)\]`).FindStringSubmatch(filename)
	if len(metaId) > 0 {
		fileInfo.MetaProvider = metaId[1]
		fileInfo.MetaId = metaId[2]
		return fileInfo
	}
	year := regexp.MustCompile(`\((\d{4})\)`).FindStringSubmatch(filenameWithoutExtension)
	if len(year) > 0 {
		yearInt, _ := strconv.Atoi(year[1])

		fileInfo.Year = yearInt

		name := regexp.MustCompile(`(.*)\(\d+\)`).FindStringSubmatch(filenameWithoutExtension)
		if len(name) > 0 {
			fileInfo.Name = name[1]
		}
	} else {
		name := regexp.MustCompile(`(.+?)(\(\d{4}\)|\[\w+-\d+\]|$)`).FindStringSubmatch(filenameWithoutExtension)
		if len(name) > 0 {
			fileInfo.Name = name[1]
		}
	}

	// Replace underscores with spaces
	fileInfo.Name = regexp.MustCompile(`_`).ReplaceAllString(fileInfo.Name, " ")
	return fileInfo
}
