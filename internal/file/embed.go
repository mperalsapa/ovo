package file

import (
	"embed"
	"io/fs"

	"log"
	"net/http"
	"os"
)

func GetFileSystem(useOS bool, staticAssets embed.FS, path string) http.FileSystem {
	if useOS {
		log.Print("Using live filesystem mode")
		return http.FS(os.DirFS(path))
	}

	log.Print("Using embed filesystem mode")
	fsys, err := fs.Sub(staticAssets, path)
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
