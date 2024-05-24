package apiController

import (
	"ovo-server/internal/model"
	"ovo-server/internal/session"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ToggleFavoriteItemRequest struct {
	ItemID uint `json:"itemID"`
}

func ToggleFavoriteItem(echo echo.Context) error {
	request := ToggleFavoriteItemRequest{}
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
