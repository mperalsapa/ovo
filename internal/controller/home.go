package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello, World!")
}
