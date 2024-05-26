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

func Library(context echo.Context) error {
	id := context.Param("id")

	order_by := context.QueryParam("order_by")
	order := context.QueryParam("order")
	var sort string
	if order_by != "" {
		sort = order_by
		if order != "" {
			sort = sort + " " + order
		}
	}

	// Getting INT from url param
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return context.Redirect(http.StatusFound, router.Routes.Home)
	}

	// Getting library data from database
	library, err := model.GetLibraryById(uint(idInt))
	if err != nil {
		log.Println(err)
		return context.Redirect(http.StatusFound, router.Routes.Home)
	}

	// Getting items from database
	library.LoadItems(sort)

	component := page.Library(page.LibraryPageData{
		Library:     library,
		UserSession: session.GetUserSession(context),
	})

	return RenderView(context, http.StatusOK, component)
}
