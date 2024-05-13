package controller

import (
	"log"
	"net/http"
	"ovo-server/internal/model"
	"ovo-server/internal/session"
	"ovo-server/internal/template/page"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Player(c echo.Context) error {
	itemID, err := strconv.Atoi(c.QueryParam("item"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid item ID")
	}
	item, err := model.GetItemById(uint(itemID))
	if err != nil {
		return c.String(http.StatusBadRequest, "Item not found")
	}

	// groupID := c.QueryParam("group")

	componentData := page.PlayerData{
		UserSession: session.GetUserSession(c),
		Item:        item,
	}

	component := page.Player(componentData)

	return RenderView(c, http.StatusOK, component)

}

func Download(c echo.Context) error {
	log.Println("Download request")
	itemID, err := strconv.Atoi(c.QueryParam("item"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid item ID")
	}
	item, err := model.GetItemById(uint(itemID))
	if err != nil {
		return c.String(http.StatusBadRequest, "Item not found")
	}

	itemAbsolutePath, err := filepath.Abs(item.FilePath)
	if err != nil {
		return c.String(http.StatusBadRequest, "Item not found")
	}

	log.Println("Download request for item: ", itemAbsolutePath)
	// r := bytes.NewReader([]byte(item.FilePath))
	// return c.Stream(http.StatusOK, "application/octet-stream", r)
	return c.File(itemAbsolutePath)
	// groupID := c.QueryParam("group")

}
