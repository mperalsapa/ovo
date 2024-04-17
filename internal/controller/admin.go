package controller

import (
	"net/http"
	"ovo-server/internal/session"
	"ovo-server/internal/template/page"

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
	// pageData := page.LibrariesPageData{}
	component := page.LibrariesPage()
	return RenderView(context, http.StatusOK, component)
}
