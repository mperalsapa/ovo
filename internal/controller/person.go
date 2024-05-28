package controller

import (
	"net/http"
	"ovo-server/internal/model"
	"ovo-server/internal/router"
	"ovo-server/internal/session"
	"ovo-server/internal/template/page"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Person(context echo.Context) error {
	personId := context.Param("id")
	// Get INT from url param
	personIdInt, err := strconv.Atoi(personId)
	if err != nil {
		return context.Redirect(http.StatusFound, router.Routes.Home)
	}

	person, err := model.GetPersonById(uint(personIdInt))
	if err != nil {
		return context.Redirect(http.StatusFound, router.Routes.Home)
	}

	person.LoadCredits()

	pageData := page.PersonDetailsPageData{
		Person:      *person,
		UserSession: session.GetUserSession(context),
	}
	component := page.PersonDetailsPage(pageData)

	return RenderView(context, http.StatusOK, component)
}
