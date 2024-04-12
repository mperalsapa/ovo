package sessions

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var Key []byte
var Store *sessions.CookieStore
var Name string

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

func IsAuth(r *http.Request) (Auth bool) {
	session, _ := Store.Get(r, Name)
	ok, auth := session.Values["authenticated"].(bool)
	fmt.Println("Checking authentication ...", auth)
	if ok {
		return auth
	}
	return false
}

func IsAdmin(r *http.Request) (Admin bool) {
	auth := IsAuth(r)
	if auth != true {
		return false
	}

	session, _ := Store.Get(r, Name)
	ok, _ := session.Values["admin"].(bool)
	return ok
}
