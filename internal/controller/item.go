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

func Item(context echo.Context) error {
	id := context.Param("id")
	// Get INT from url param
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return context.Redirect(http.StatusFound, router.Routes.Home)
	}

	// Getting Item data from database
	item, err := model.GetItemById(uint(idInt))
	if err != nil {
		return context.Redirect(http.StatusFound, router.Routes.Home)
	}

	pageData := page.ItemDetailsPageData{
		Item:        item,
		UserSession: session.GetUserSession(context),
	}

	component := page.ItemDetailsPage(pageData)

	return RenderView(context, http.StatusOK, component)
}
