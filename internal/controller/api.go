package controller

import (
	"log"
	"net/http"
	"ovo-server/internal/model"

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

	err := library.SaveLibrary()
	if err != nil {
		return echo.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return echo.JSON(http.StatusOK, library)
}
