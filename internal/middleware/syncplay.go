package middleware

import (
	"ovo-server/internal/session"
	"ovo-server/internal/syncplay"

	"github.com/labstack/echo/v4"
)

func ValidateCurrentSyncplayGroup(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if user is in a sync group
		userSession := session.GetUserSession(c)
		if userSession.SyncPlayGroup != "" {
			// if user is within a group, we check if that group exists
			group := syncplay.Groups.GetGroup(userSession.SyncPlayGroup)
			if group == nil {
				// If group does not exist, we remove the user from the group
				userSession.SyncPlayGroup = ""
				userSession.SaveUserSession(c)
			}
		}
		return next(c)
	}
}
