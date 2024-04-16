package session

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var Key []byte
var Store *sessions.CookieStore
var Name string

type UserSession struct {
	Username      string
	Authenticated bool
	Role          string
	ErrorMsg      string
}

func GenerateSessionHandler(key string, name string) {
	Key = []byte(key)
	Store = sessions.NewCookieStore(Key)
	Name = name
}

func GetUsername(r *http.Request) (UserName string) {
	session, _ := Store.Get(r, Name)
	if session.IsNew {
		return ""
	}
	UserName = session.Values["username"].(string)
	return UserName
}

func GetUserSession(c echo.Context) (User UserSession) {
	session, _ := Store.Get(c.Request(), Name)
	if session.IsNew {
		return UserSession{}
	}
	user := UserSession{
		Username:      session.Values["username"].(string),
		Authenticated: session.Values["authenticated"].(bool),
		Role:          session.Values["role"].(string),
		ErrorMsg:      session.Values["error_msg"].(string),
	}
	return user
}

func (u *UserSession) SaveUserSession(c echo.Context) {
	fmt.Println("Saving user session ...")
	fmt.Println("Username : ", u.Username)
	session, _ := Store.Get(c.Request(), Name)
	session.Values["username"] = u.Username
	session.Values["authenticated"] = u.Authenticated
	session.Values["role"] = u.Role
	session.Values["error_msg"] = u.ErrorMsg
	session.Save(c.Request(), c.Response())
}

func (u *UserSession) ClearUserSession(r *http.Request, w http.ResponseWriter) {
	session, _ := Store.Get(r, Name)
	session.Values["username"] = ""
	session.Values["authenticated"] = false
	session.Values["role"] = ""
	session.Values["error_msg"] = ""
	session.Save(r, w)
}

func GetKey(c echo.Context) (Key string) {
	session, _ := Store.Get(c.Request(), Name)
	if session.IsNew {
		return ""
	}
	Key = session.Values["key"].(string)
	return Key
}

func SetKey(c echo.Context, key string, value string) {
	session, _ := Store.Get(c.Request(), Name)
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
	fmt.Println("Checking authentication ...", auth)
	return auth
}

func IsAdmin(c echo.Context) (Admin bool) {
	auth := IsAuth(c)
	if auth != true {
		return false
	}

	userSession := GetUserSession(c)
	role := userSession.Role
	if role != "admin" {
		return false
	}

	return true
}
