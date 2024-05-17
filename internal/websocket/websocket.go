package websocket

import (
	"errors"
	"log"
	"net/http"
	"os"
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

func (s *WsServer) AddClient(conn *websocket.Conn) {
	s.Lock()
	s.clients[conn] = true
	s.Unlock()
	log.Println("Added new client. Total clients: ", len(s.clients))
}

func (s *WsServer) RemoveClient(conn *websocket.Conn) {
	s.Lock()
	if _, ok := s.clients[conn]; ok {
		delete(s.clients, conn)
		conn.Close()
	}
	s.Unlock()
	log.Println("Removed client. Total clients: ", len(s.clients))
}

func (s *WsServer) ReadLoop(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {

			err := ws.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				// La conexión está cerrada
				log.Println("La conexión WebSocket está cerrada.")
				break
			}

			log.Println("Error reading message: ", err)
			continue
		}

		log.Println("Received message: ", string(message))
		response := []byte("Received message: " + string(message))
		ws.WriteMessage(websocket.TextMessage, response)
	}
}

func (s *WsServer) Broadcast(message []byte) {
	for conn := range s.clients {
		log.Printf("Broadcasting message to %s.\n", conn.RemoteAddr().String())
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error broadcasting message: ", err)
			conn.Close()
			delete(s.clients, conn)
		}
	}
}

func BroadcastToGroup(sender *websocket.Conn, group []*websocket.Conn, message []byte) {
	for _, conn := range group {
		if conn == sender {
			continue
		}

		log.Printf("Broadcasting message to %s.\n", conn.RemoteAddr().String())
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error broadcasting message: ", err)
			conn.Close()
		}
	}
}
