package syncplay

import (
	"log"
	"ovo-server/internal/model"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Groups SyncGroups

type Sync struct {
	CurrentItem *model.Item
	CurrentTime int
	StartedFrom int
	IsPlaying   bool
}

type SyncGroup struct {
	ID          uuid.UUID `json:"id"`
	Name        string
	Connections []*websocket.Conn
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
		Connections: make([]*websocket.Conn, 0),
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

func (sg *SyncGroups) AddConnectionToGroup(id string, conn *websocket.Conn) {
	sg.Lock()
	defer sg.Unlock()

	group := sg.GetGroup(id)
	group.Connections = append(group.Connections, conn)
}

func (sg *SyncGroups) RemoveConnectionFromGroup(id string, conn *websocket.Conn) {
	sg.Lock()
	defer sg.Unlock()

	group := sg.GetGroup(id)
	connIndex := -1

	for i, c := range group.Connections {
		if c == conn {
			connIndex = i
			break
		}
	}

	if connIndex == -1 {
		return
	}

	group.Connections = append(group.Connections[:connIndex], group.Connections[connIndex+1])
}

func (s *Sync) SetNewItem(item *model.Item) {
	if s.CurrentItem != item {
		s.CurrentItem = item
		s.CurrentTime = 0
		s.StartedFrom = 0
		s.IsPlaying = false
	}
}
