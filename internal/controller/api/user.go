package apiController

import (
	"ovo-server/internal/model"
	"ovo-server/internal/session"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ToggleItemUserState struct {
	ItemID uint `json:"itemID"`
}

func ToggleFavoriteItem(echo echo.Context) error {
	request := ToggleItemUserState{}
	err := echo.Bind(&request)

	if err != nil {
		return echo.JSON(400, map[string]string{"error": "Invalid payload"})
	}

	userSession := session.GetUserSession(echo)
	user := model.GetUserByUsername(userSession.Username)
	item, err := model.GetItemById(request.ItemID)

	if err != nil {
		return echo.JSON(400, map[string]string{"error": "Item not found"})
	}

	isFavorite := user.ToggleFavoriteItem(item.ID)

	return echo.JSON(200, map[string]string{"message": "success", "isFavorite": strconv.FormatBool(isFavorite)})
}

func ToggleWatchedItem(echo echo.Context) error {
	request := ToggleItemUserState{}
	err := echo.Bind(&request)

	if err != nil {
		return echo.JSON(400, map[string]string{"error": "Invalid payload"})
	}

	userSession := session.GetUserSession(echo)
	user := model.GetUserByUsername(userSession.Username)
	item, err := model.GetItemById(request.ItemID)

	if err != nil {
		return echo.JSON(400, map[string]string{"error": "Item not found"})
	}

	isWatched := user.ToggleWatchedItem(item.ID)

	return echo.JSON(200, map[string]string{"message": "success", "isWatched": strconv.FormatBool(isWatched)})
}
