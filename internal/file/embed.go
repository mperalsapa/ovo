package file

import (
	"embed"
	"io/fs"

	"log"
	"net/http"
	"os"
)

func GetFileSystem(useOS bool, staticAssets embed.FS) http.FileSystem {
	if useOS {
		log.Print("Using live filesystem mode")
		return http.FS(os.DirFS("public"))
	}

	log.Print("Using embed filesystem mode")
	fsys, err := fs.Sub(staticAssets, "public")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
