package controller

import (
	"fmt"
	"net/http"
	"ovo-server/internal/model"
	session "ovo-server/internal/session"

	"github.com/labstack/echo/v4"
)

func Login(e echo.Context) error {
	var reqUser model.User
	if err := e.Bind(&reqUser); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}
	fmt.Println("Checking against : ", reqUser.Username)
	user := model.GetUserByUsername(reqUser.Username)

	fmt.Println("User : \n", user, "ReqUser : \n", reqUser)
	if valid := user.CheckPassword(reqUser.Password); !valid {
		return e.JSON(http.StatusUnauthorized, "Invalid username or password")
	}

	session, _ := session.Store.Get(e.Request(), session.Name)
	session.Values["authenticated"] = true
	session.Values["username"] = user.Username

	session.Save(e.Request(), e.Response())

	return e.JSON(http.StatusOK, user)
}

func Register(e echo.Context) error {
	var reqUser model.User
	if err := e.Bind(&reqUser); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	user := model.CreateUser(reqUser.Username, reqUser.Password, reqUser.Role)
	user.Save()

	return e.JSON(http.StatusOK, user)
}

func Logout(e echo.Context) error {
	session, _ := session.Store.Get(e.Request(), session.Name)
	session.Values["authenticated"] = false
	session.Save(e.Request(), e.Response())

	fmt.Println(session.Values["username"], " Logged out ", session.Values["authenticated"])
	return e.JSON(http.StatusOK, "Logged out")
}

func About(w http.ResponseWriter, r *http.Request) {
	msg := "Hello from About, World!"
	w.Write([]byte(msg))
}
