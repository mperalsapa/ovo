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
		log.Print("using live mode")
		return http.FS(os.DirFS("public"))
	}

	log.Print("using embed mode")
	fsys, err := fs.Sub(staticAssets, "public")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
