package controller

import (
	"net/http"
	"ovo-server/internal/session"
	"ovo-server/internal/template/page"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	userSession := session.GetUserSession(c)

	pageData := page.HomePageData{
		UserSession: userSession,
	}
	component := page.HomePage(pageData)
	return RenderView(c, http.StatusOK, component)
}
