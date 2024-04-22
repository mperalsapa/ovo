package file

import (
	"io/ioutil"
	"log"
)

type FileInfo struct {
	Name string
	Path string
	Size int64
}

// This function scans a dir and returns a slice of files
func ScanPath(path string) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		log.Println(file.Name())
	}

}
