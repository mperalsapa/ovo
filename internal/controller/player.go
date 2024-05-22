package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"ovo-server/internal/model"
	"ovo-server/internal/session"
	"ovo-server/internal/syncplay"
	"ovo-server/internal/template/page"
	"ovo-server/internal/websocket"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PlayerData struct {
	ItemID uint `query:"item"`
}

func Player(c echo.Context) error {
	var playerData PlayerData
	var item model.Item

	userSession := session.GetUserSession(c)
	c.Bind(&playerData)
	itemID := playerData.ItemID

	// If item id is provided in the query, we load the item from database
	if itemID != 0 {
		item, _ = model.GetItemById(itemID)
	}

	// After checking for new item, we check if user is in a sync group
	if userSession.SyncPlayGroup != "" {
		group := syncplay.Groups.GetGroup(userSession.SyncPlayGroup)

		// If user in group wants to play a new item, we sync the item
		if item.ID != 0 {
			group.Sync.SetNewItem(&item)
			message := websocket.Message{
				Event: "newItem",
				Item:  &item,
			}
			messageData, _ := json.Marshal(message)
			websocket.BroadcastToList(group.Connections, messageData, nil)
		}

		if group.Sync.CurrentItem != nil {
			item = *group.Sync.CurrentItem
		}

	}

	componentData := page.PlayerData{
		UserSession: session.GetUserSession(c),
		Item:        item,
	}

	component := page.Player(componentData)

	return RenderView(c, http.StatusOK, component)

}

func Download(c echo.Context) error {
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
	return c.File(itemAbsolutePath)

}
