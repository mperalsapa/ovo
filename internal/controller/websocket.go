package controller

import (
	"log"

	"ovo-server/internal/websocket"

	"github.com/labstack/echo/v4"
)

func WebsocketHandler(context echo.Context) error {
	ws, err := websocket.Server.Upgrade(context.Response(), context.Request())
	if err != nil {
		return err
	}
	log.Println("New websocket connection")
	websocket.Server.AddClient(ws)
	websocket.Server.ReadLoop(ws)
	defer websocket.Server.RemoveClient(ws)
	return nil
}
