package session

import (
	"encoding/gob"
	"fmt"
	"log"
	"ovo-server/internal/model"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var SessionSettings sessionSettings

type sessionSettings struct {
	Store  *sessions.CookieStore
	Name   string `default:"ovo-session"`
	MaxAge int    `default:"3600"`
}

type UserSession struct {
	Username      string
	Authenticated bool
	SyncPlayGroup string
	Role          model.Role
	ErrorMsg      string
	DeviceID      uuid.UUID
}

func GenerateSessionHandler(key string, name string) {

	SessionSettings = sessionSettings{
		Store:  sessions.NewCookieStore([]byte(key)),
		Name:   name,
		MaxAge: 3600,
	}

	gob.Register(UserSession{})

}

func GetUserSession(c echo.Context) (User UserSession) {
	currentSession, err := SessionSettings.Store.Get(c.Request(), SessionSettings.Name)

	if err != nil {
		log.Println("Error getting session: ", err)
		return UserSession{}
	}

	if currentSession.IsNew {
		return UserSession{}
	}
	user, ok := currentSession.Values["user"].(UserSession)
	if !ok {
		log.Println("Error getting user from session")
		return UserSession{}
	}

	return user
}

func GetSession(c echo.Context) (*sessions.Session, error) {
	return SessionSettings.Store.Get(c.Request(), SessionSettings.Name)
}

func (u *UserSession) SaveUserSession(c echo.Context) {

	session, err := GetSession(c)
	if err != nil {
		log.Printf("Error getting session when storing: %s, overwriting session with a new one.", err)
		// if we could not get the session, we overwrite the cookie
		session = sessions.NewSession(SessionSettings.Store, SessionSettings.Name)
	}

	log.Printf("Storing user %s with role %d is authenticated: %t", u.Username, u.Role, u.Authenticated)
	session.Values["user"] = *u
	err = session.Save(c.Request(), c.Response())
	if err != nil {
		log.Println("Error saving session: ", err)
	}
	log.Println("User session saved: " + u.Username + " - Authenticated: " + fmt.Sprint(u.Authenticated))
}

func GetKey(c echo.Context) (Key string) {
	session, err := GetSession(c)
	if err != nil {
		log.Println("Error getting session when getting key: ", err)
		return ""
	}

	if session.IsNew {
		return ""
	}
	Key = session.Values["key"].(string)
	return Key
}

func SetKey(c echo.Context, key string, value string) {
	session, err := GetSession(c)
	if err != nil {
		log.Println("Error getting session when setting key: ", err)
		return
	}

	session.Values["key"] = key
	session.Save(c.Request(), nil)
}

func (u *UserSession) PopErrorMessage(c echo.Context) string {
	errorMsg := u.ErrorMsg
	u.ErrorMsg = ""
	u.SaveUserSession(c)
	return errorMsg
}

func IsAuth(c echo.Context) (Auth bool) {
	userSession := GetUserSession(c)
	auth := userSession.Authenticated
	return auth
}

func IsAdmin(c echo.Context) (Admin bool) {
	auth := IsAuth(c)
	if !auth {
		return false
	}

	userSession := GetUserSession(c)
	role := userSession.Role
	return role == model.Admin
}

func (userSession *UserSession) IsAdmin() bool {
	return userSession.Role == model.Admin
}
