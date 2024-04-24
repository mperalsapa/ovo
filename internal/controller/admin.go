package controller

import (
	"log"
	"net/http"
	"ovo-server/internal/model"
	"ovo-server/internal/router"
	"ovo-server/internal/session"
	"ovo-server/internal/template/page"
	"strconv"

	"github.com/labstack/echo/v4"
)

func AdminDashboard(context echo.Context) error {
	pageData := page.AdminDashboardPageData{
		Username: session.GetUserSession(context).Username,
	}
	component := page.AdminDashboardPage(pageData)
	return RenderView(context, http.StatusOK, component)
}

func AdminLibraries(context echo.Context) error {
	pageData := page.LibrariesPageData{}
	pageData.Libraries = model.GetLibraries()
	component := page.LibrariesPage(pageData)
	return RenderView(context, http.StatusOK, component)
}

func AdminLibrary(context echo.Context) error {
	idStr := context.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return context.Redirect(http.StatusFound, router.AdminRoutes.Libraries)
	}

	library := model.Library{}
	library, err = model.GetLibraryById(uint(id))

	if err != nil && id != 0 {
		return context.Redirect(http.StatusFound, router.AdminRoutes.Libraries)
	}

	pageData := page.AdminLibraryFormPageData{
		Library: library,
		Editing: library.ID != 0,
	}

	component := page.AdminLibraryForm(pageData)
	return RenderView(context, http.StatusOK, component)
}

type LibraryForm struct {
	Submit string   `form:"submit"`
	ID     uint     `form:"id"`
	Type   string   `form:"type"`
	Name   string   `form:"name"`
	Paths  []string `form:"paths[]"`
}

func AdminStoreLibrary(context echo.Context) error {

	libraryForm := LibraryForm{}

	if err := context.Bind(&libraryForm); err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	// query library. In case of error and ID is not 0, redirect to libraries (meaning the library does not exist)
	// otherwise create a new library
	storedLibrary, err := model.GetLibraryById(libraryForm.ID)
	if err != nil && libraryForm.ID != 0 {
		return context.Redirect(http.StatusFound, router.AdminRoutes.Libraries)
	}

	// if we are deleting the library, delete it and redirect to libraries
	if libraryForm.Submit == "Delete" {
		storedLibrary.DeleteLibrary()
		return context.Redirect(http.StatusFound, router.AdminRoutes.Libraries)
	}

	// Storing library data into the storedLibrary variable
	switch libraryForm.Type {
	case string(model.LibraryTypeMovie):
		storedLibrary.Type = model.LibraryTypeMovie
	case string(model.LibraryTypeShow):
		storedLibrary.Type = model.LibraryTypeShow
	}
	storedLibrary.Name = libraryForm.Name
	storedLibrary.Paths = libraryForm.Paths

	// Trying to save the library
	err = storedLibrary.SaveLibrary()
	// if there is an error, show the form again with the error message and the previous data
	if err != nil {
		pageData := page.AdminLibraryFormPageData{
			Library:  storedLibrary,
			ErrorMsg: err.Error(),
			Editing:  storedLibrary.ID != 0,
		}
		component := page.AdminLibraryForm(pageData)
		return RenderView(context, http.StatusBadRequest, component)
	}

	// if everything is ok, redirect to libraries
	return context.Redirect(http.StatusFound, router.AdminRoutes.Libraries)
}

func AdminCommand(context echo.Context) error {
	command := context.Param("action")
	log.Println("Received command request: ", command)

	switch command {
	case "ScanLibraries":
		log.Println("Scanning libraries")
		libraries := model.GetLibraries()
		for _, library := range libraries {
			log.Println("Scanning library: ", library.Name)
			library.ScanLibrary()
			// for _, path := range library.Paths {
			// 	files := file.ScanPath(path)
			// 	for _, f := range files {
			// 		log.Println("Found file: ", f)
			// 		fileInfo := file.ParseFilename(f)
			// 		log.Printf("Name: %s, Year: %d, Meta Prov.: %s, Meta ID: %s", fileInfo.Name, fileInfo.Year, fileInfo.MetaProvider, fileInfo.MetaId)
			// 		movie := tmdb.FindMovieByFileInfo(fileInfo)

			// 		if movie != nil {
			// 			log.Printf("Found movie: %s (%s)\n%s\n\n", movie.Title, movie.ReleaseDate, movie.Overview)
			// 		} else {

			// 		}
			// 	}
			// }
		}
	default:
		log.Println("Unknown command")
	}

	return context.Redirect(http.StatusFound, router.AdminRoutes.Dashboard)
}
