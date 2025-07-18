package controller

import (
	"log"
	"net/http"
	"ovo-server/internal/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

func APIGetLibraries(echo echo.Context) error {
	libraries := model.GetLibraries()
	return echo.JSON(http.StatusOK, libraries)
}

func APIAddLibrary(echo echo.Context) error {
	library := model.Library{}
	log.Println(echo.Request().Body)
	if err := echo.Bind(&library); err != nil {
		log.Println(err)
		return echo.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	err := library.Save()
	if err != nil {
		return echo.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return echo.JSON(http.StatusOK, library)
}

func APIDeleteLibrary(echo echo.Context) error {
	idStr := echo.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	err = model.DeleteLibrary(uint(id))
	if err != nil {
		return echo.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return echo.JSON(http.StatusOK, map[string]string{"message": "Library deleted"})
}

func APIGetLibrary(echo echo.Context) error {
	idStr := echo.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	library, err := model.GetLibraryById(uint(id))
	if err != nil && id != 0 {
		return echo.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return echo.JSON(http.StatusOK, library)
}
