package controller

import (
	"log"

	"ovo-server/internal/session"
	"ovo-server/internal/syncplay"
	"ovo-server/internal/websocket"

	"github.com/labstack/echo/v4"
)

func WebsocketHandler(context echo.Context) error {
	// First we check if the user is within a sync group
	usersession := session.GetUserSession(context)
	if usersession.SyncPlayGroup == "" {
		return echo.NewHTTPError(400, "User not in a sync group")
	}

	group := syncplay.Groups.GetGroup(usersession.SyncPlayGroup)
	if group == nil {
		return echo.NewHTTPError(400, "User not in a sync group")
	}

	ws, err := websocket.Server.Upgrade(context.Response(), context.Request())
	if err != nil {
		return err
	}
	log.Println("New websocket connection")

	group.AddConnection(ws)
	websocket.Server.ReadLoop(ws, group)
	defer group.RemoveConnection(ws)
	return nil
}
