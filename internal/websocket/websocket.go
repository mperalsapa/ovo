package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"ovo-server/internal/model"
	"ovo-server/internal/syncplay"
	"sync"

	"github.com/gorilla/websocket"
)

var Server *WsServer

type WsServer struct {
	sync.Mutex
	clients  map[*websocket.Conn]bool
	upgrader *websocket.Upgrader
}

func Init() {
	Server = NewWsServer()
	if Server != nil {
		log.Println("Websocket server initialized.")
		return
	}

	// websocket server could not start. Stopping the server.
	log.Println("Error initializing websocket server.")
	os.Exit(1)
}

func NewWsServer() *WsServer {
	server := new(WsServer)
	server.clients = make(map[*websocket.Conn]bool)
	server.upgrader = new(websocket.Upgrader)
	return server
}

func (s *WsServer) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	if s == nil {
		return nil, errors.New("websocket server is nil")
	}

	if w == nil || r == nil {
		return nil, errors.New("responseWriter or request is nil")
	}

	conn, err := Server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
		return nil, err
	}

	return conn, nil
}

type Message struct {
	Event       string      `json:"event"`
	StartedFrom float32     `json:"StartedFrom"`
	StartedAt   int64       `json:"StartedAt,omitempty"`
	Item        *model.Item `json:"Item,omitempty"`
}

func (s *WsServer) ReadLoop(ws *websocket.Conn, group *syncplay.SyncGroup) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {

			err := ws.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				// Websocket connection is closed
				log.Println("Websocket connection closed.")
				break
			}

			log.Println("Error reading message: ", err)
			continue
		}

		var parsedMessage Message
		err = json.Unmarshal(message, &parsedMessage)
		if err != nil {
			log.Println("Error unmarshalling message: ", err)
			continue
		}
		log.Println("Received message: ", parsedMessage.Event)
		log.Printf("Waiting connections: %d", group.GetWaitingCount())
		switch parsedMessage.Event {
		case "ping":
			SendMessage(ws, Message{Event: "pong"})
			continue
		case "requestPlay":
			var action string
			if group.Sync.IsPlaying && group.CanPlay() {
				action = "play"
			} else {
				action = "pause"
			}
			SendMessage(ws, Message{Event: action, StartedFrom: group.Sync.StartedFrom, StartedAt: group.Sync.StartedAt})
			continue
		case "play":
			group.SetState(true, parsedMessage.StartedFrom)
			group.SetPlayingBeforeBuffer(true)
			parsedMessage.StartedFrom = group.Sync.StartedFrom
			parsedMessage.StartedAt = group.Sync.StartedAt
		case "pause":
			group.SetState(false, parsedMessage.StartedFrom)
			group.SetPlayingBeforeBuffer(false)
			parsedMessage.StartedFrom = group.Sync.StartedFrom
		case "seek":
			group.SetState(true, parsedMessage.StartedFrom)
			parsedMessage.StartedFrom = group.Sync.StartedFrom
			parsedMessage.StartedAt = group.Sync.StartedAt
		case "buffering":
			group.SetState(false, parsedMessage.StartedFrom)
			group.SetWaitingClient(ws)
			parsedMessage.StartedFrom = group.Sync.StartedFrom
		case "canplay":
			group.SetReadyClient(ws)
			if !group.CanPlay() {
				continue
			}

			log.Println("Group can play")
			group.SetState(true, group.Sync.StartedFrom)
			parsedMessage.StartedFrom = group.Sync.StartedFrom
			parsedMessage.StartedAt = group.Sync.StartedAt
			if group.Sync.IsPlayingBeforeBuffer {
				parsedMessage.Event = "play"
			}
		}

		stringifiedMessage, err := json.Marshal(parsedMessage)
		if err != nil {
			log.Println("Error marshalling message: ", err)
			continue
		}

		BroadcastToList(group.GetConnectionList(), stringifiedMessage, nil)
	}
}

func SendMessage(ws *websocket.Conn, message Message) {
	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message: ", err)
		return
	}

	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Println("Error sending message: ", err)
		ws.Close()
	}
}

func BroadcastToList(connections []*websocket.Conn, message []byte, sender *websocket.Conn) {
	for _, conn := range connections {
		if sender != nil && conn == sender {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error broadcasting message: ", err)
			conn.Close()
		}
	}
}
