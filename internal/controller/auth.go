package controller

import (
	"fmt"
	"net/http"
	"ovo-server/internal/model"
	"ovo-server/internal/router"
	session "ovo-server/internal/session"
	"ovo-server/internal/view"

	"github.com/labstack/echo/v4"
)

func Login(e echo.Context) error {
	userSession := session.GetUserSession(e)
	fmt.Println("Username stored in session : ", userSession.Username)
	component := view.LoginPage(userSession.Username)
	return RenderView(e, http.StatusOK, component)
}

func LoginRequest(e echo.Context) error {
	var reqUser model.User
	if err := e.Bind(&reqUser); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	userSession := session.GetUserSession(e)
	userSession.Username = reqUser.Username

	user := model.GetUserByUsername(reqUser.Username)

	if valid := user.CheckPassword(reqUser.Password); !valid {
		userSession.Authenticated = false
		userSession.SaveUserSession(e)
		return e.JSON(http.StatusUnauthorized, "Invalid username or password")
	}
	userSession.Authenticated = true
	userSession.SaveUserSession(e)

	return e.Redirect(http.StatusFound, router.Routes.Home)
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

func Logout(c echo.Context) error {
	session := session.GetUserSession(c)
	session.Authenticated = false
	session.SaveUserSession(c)
	fmt.Println("User logged out: " + session.Username)

	return c.Redirect(http.StatusFound, router.Routes.Login)
}

func About(w http.ResponseWriter, r *http.Request) {
	msg := "Hello from About, World!"
	w.Write([]byte(msg))
}
