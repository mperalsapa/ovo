package controller

import (
	"net/http"
	"ovo-server/internal/template/page"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	pageData := page.HomePageData{
		Username: "user",
	}
	component := page.HomePage(pageData)
	return RenderView(c, http.StatusOK, component)
}
