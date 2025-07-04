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
	MetaID       string
	FilePath     string
}

// This function scans a dir and returns a slice of files
func ScanFiles(path string) []string {
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

// This function scans the given path and returns a slice of Directories
func ScanDirectories(path string) []string {
	entries, err := os.ReadDir(path)
	dirs := []string{}
	if err != nil {
		log.Println(err)
	}

	for _, file := range entries {
		if !file.IsDir() {
			continue
		}

		dirs = append(dirs, file.Name())
	}

	return dirs
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
	fileInfo.FilePath = filename

	// removing file extension
	filenameWithoutExtension := regexp.MustCompile(`(.+?)(\.[^.]*$|$)`).FindStringSubmatch(filename)[1]

	metaId := regexp.MustCompile(`\[(.+)id-(.+)\]`).FindStringSubmatch(filename)
	if len(metaId) > 0 {
		fileInfo.MetaProvider = metaId[1]
		fileInfo.MetaID = metaId[2]
		return fileInfo
	}

	year := ParseYearFromFilename(filename)
	if year != 0 {

		fileInfo.Year = year

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

// This function tries to get the year from the filename
// It expects the year to be in parentheses
// Example: "The Matrix (1999)"
func ParseYearFromFilename(filename string) int {
	year := regexp.MustCompile(`\((\d{4})\)`).FindStringSubmatch(filename)
	if len(year) > 0 {
		yearInt, _ := strconv.Atoi(year[1])
		return yearInt
	}

	return 0
}

func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}
