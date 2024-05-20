package apiController

import (
	"log"
	"ovo-server/internal/model"
	"ovo-server/internal/session"
	"ovo-server/internal/syncplay"

	"github.com/labstack/echo/v4"
)

type GetSyncGroupsResponse struct {
	Groups       map[string]*syncplay.SyncGroup `json:"groups"`
	CurrentGroup string                         `json:"currentGroup"`
}

func GetSyncGroups(context echo.Context) error {
	response := GetSyncGroupsResponse{}
	response.Groups = syncplay.Groups.Groups

	userSession := session.GetUserSession(context)
	if userSession.SyncPlayGroup != "" && syncplay.Groups.GetGroup(userSession.SyncPlayGroup) != nil {
		response.CurrentGroup = userSession.SyncPlayGroup
	}

	return context.JSON(200, response)
}

func CreateSyncGroup(context echo.Context) error {
	userSession := session.GetUserSession(context)

	// First, we check if user is already in a group
	group := syncplay.Groups.GetGroup(userSession.SyncPlayGroup)

	if group != nil {
		return context.JSON(400, map[string]string{"error": "User is already in a group"})
	}

	user := model.GetUserByUsername(userSession.Username)
	group = syncplay.Groups.CreateGroup(user)
	userSession.SyncPlayGroup = group.ID.String()
	userSession.SaveUserSession(context)

	return context.JSON(200, map[string]string{"message": "Group created"})
}

type SyncGroupJoin struct {
	GroupID string `json:"id"`
}

func JoinSyncGroup(context echo.Context) error {
	userSession := session.GetUserSession(context)
	syncGroup := SyncGroupJoin{}
	err := context.Bind(&syncGroup)
	if err != nil {
		return context.JSON(400, map[string]string{"error": "Invalid JSON"})
	}
	log.Println(syncGroup)
	if syncGroup.GroupID == "" {
		return context.JSON(400, map[string]string{"error": "Group ID is required"})
	}

	group := syncplay.Groups.GetGroup(syncGroup.GroupID)
	if group == nil {
		return context.JSON(400, map[string]string{"error": "Group not found"})
	}

	group.AddUser(userSession.Username)
	userSession.SyncPlayGroup = group.ID.String()
	userSession.SaveUserSession(context)
	return context.JSON(200, map[string]string{"message": "Group joined"})

}

func LeaveSyncGroup(context echo.Context) error {
	userSession := session.GetUserSession(context)

	group := syncplay.Groups.GetGroup(userSession.SyncPlayGroup)

	userSession.SyncPlayGroup = ""
	userSession.SaveUserSession(context)

	if group != nil {
		group.RemoveUser(userSession.Username)
	}
	return context.JSON(200, map[string]string{"message": "User left group"})
}
