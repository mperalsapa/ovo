package syncplay

import (
	"log"
	"ovo-server/internal/model"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Groups SyncGroups

type Sync struct {
	sync.RWMutex
	CurrentItem           *model.Item
	StartedFrom           float32 // Time in seconds from the start of the video. If the video started at half way and its complete duration is 100s, would be 50s.
	StartedAt             int64   // Unix (miliseconds) timestamp when the video started playing. Adding startedFrom to this would give the current time.
	IsPlaying             bool    `default:"true"`
	IsPlayingBeforeBuffer bool    `default:"true"`
}

type SyncGroup struct {
	sync.RWMutex
	ID          uuid.UUID `json:"id"`
	Name        string
	Connections map[*websocket.Conn]bool `json:"-"` // Map of connections in the group, and their ready status (Ready to play, or not).
	Users       []string
	Sync        Sync
}

type SyncGroups struct {
	Groups map[string]*SyncGroup
	sync.RWMutex
}

func Init() {
	Groups = SyncGroups{
		Groups: make(map[string]*SyncGroup),
	}
}

func (sg *SyncGroups) CreateGroup(user model.User) *SyncGroup {
	newID, err := uuid.NewV7()

	if err != nil {
		log.Println("Error creating new group ID: ", err)
		return nil
	}
	checkGroup := sg.GetGroup(newID.String())

	if checkGroup != nil {
		return nil
	}
	sg.Lock()
	defer sg.Unlock()

	group := &SyncGroup{
		ID:          newID,
		Name:        user.Username + "'s group",
		Connections: make(map[*websocket.Conn]bool, 0),
		Users:       []string{user.Username},
		Sync:        Sync{},
	}
	sg.Groups[newID.String()] = group

	return group

}

func (sg *SyncGroups) GetGroup(id string) *SyncGroup {
	sg.RLock()
	defer sg.RUnlock()

	return sg.Groups[id]
}

func (sg *SyncGroups) GetGroups() map[string]*SyncGroup {
	sg.RLock()
	defer sg.RUnlock()

	return sg.Groups
}

func (sg *SyncGroups) DeleteGroup(id string) {
	sg.Lock()
	defer sg.Unlock()

	delete(sg.Groups, id)
}

func (g *SyncGroup) AddUser(username string) {
	g.Users = append(g.Users, username)
}

func (g *SyncGroup) RemoveUser(user string) {
	userIndex := -1

	for i, u := range g.Users {
		if u == user {
			userIndex = i
			break
		}
	}

	if userIndex == -1 {
		return
	}

	g.Users = append(g.Users[:userIndex], g.Users[userIndex+1:]...)

	if len(g.Users) == 0 {
		Groups.DeleteGroup(g.ID.String())
	}
}

func (g *SyncGroup) AddConnection(conn *websocket.Conn) {
	g.Lock()
	defer g.Unlock()
	g.Sync.Lock()
	defer g.Sync.Unlock()

	if len(g.Connections) == 0 && g.Sync.CurrentItem != nil {
		g.Sync.StartedAt = time.Now().UnixMilli()
		g.Sync.StartedFrom = 0
	}

	g.Connections[conn] = true
}

func (g *SyncGroup) RemoveConnection(conn *websocket.Conn) {
	g.Lock()
	defer g.Unlock()

	delete(g.Connections, conn)
}

func (g *SyncGroup) GetConnectionList() []*websocket.Conn {
	var connections []*websocket.Conn
	for k := range g.Connections {
		connections = append(connections, k)
	}
	return connections
}

func (g *SyncGroup) SetWaitingClient(conn *websocket.Conn) {
	g.Lock()
	defer g.Unlock()

	g.Connections[conn] = false
}

func (g *SyncGroup) SetReadyClient(conn *websocket.Conn) {
	g.Lock()
	defer g.Unlock()

	g.Connections[conn] = true
}

func (g *SyncGroup) GetWaitingCount() int {
	g.RLock()
	defer g.RUnlock()

	count := 0
	for _, ready := range g.Connections {
		if !ready {
			count++
		}
	}

	return count

}

func (g *SyncGroup) CanPlay() bool {
	return g.GetWaitingCount() == 0
}

func (g *SyncGroup) SetState(playing bool, runtime float32) {
	canplay := g.CanPlay()
	g.Lock()
	defer g.Unlock()

	if canplay {
		g.Sync.Lock()
		defer g.Sync.Unlock()

		g.Sync.IsPlaying = true
		g.Sync.StartedFrom = runtime
		g.Sync.StartedAt = time.Now().UnixMilli()
	}
}

func (g *SyncGroup) SetPlayingBeforeBuffer(playingBeforeBuffer bool) {
	g.Lock()
	defer g.Unlock()

	g.Sync.Lock()
	defer g.Sync.Unlock()

	g.Sync.IsPlayingBeforeBuffer = playingBeforeBuffer
}

func (s *Sync) SetNewItem(item *model.Item) {
	s.Lock()
	defer s.Unlock()

	if s.CurrentItem != item {
		s.CurrentItem = item
		s.StartedAt = 0
		s.StartedFrom = 0
		s.IsPlaying = true
	}
}

func (s *Sync) GetStartedAt() int64 {
	s.RLock()
	defer s.RUnlock()

	if s.IsPlaying {
		return s.StartedAt
	} else {
		return time.Now().UnixMilli()
	}
}
