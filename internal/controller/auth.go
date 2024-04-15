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

func Login(context echo.Context) error {
	userSession := session.GetUserSession(context)
	fmt.Println("Username stored in session : ", userSession.Username)
	component := view.LoginPage(userSession.Username)
	return RenderView(context, http.StatusOK, component)
}

func LoginRequest(context echo.Context) error {
	var reqUser model.User
	if err := context.Bind(&reqUser); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	userSession := session.GetUserSession(context)
	userSession.Username = reqUser.Username

	user := model.GetUserByUsername(reqUser.Username)

	if valid := user.CheckPassword(reqUser.Password); !valid {
		userSession.Authenticated = false
		userSession.ErrorMsg = "Invalid username or password"
		userSession.SaveUserSession(context)
		return context.Redirect(http.StatusFound, router.Routes.Login)
	}
	userSession.Authenticated = true
	userSession.SaveUserSession(context)

	return context.Redirect(http.StatusFound, router.Routes.Home)
}

func Register(context echo.Context) error {
	var reqUser model.User
	if err := context.Bind(&reqUser); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	user := model.CreateUser(reqUser.Username, reqUser.Password, reqUser.Role)
	user.Save()

	return context.Redirect(http.StatusFound, router.Routes.Login)
}

func Logout(context echo.Context) error {
	session := session.GetUserSession(context)
	session.Authenticated = false
	session.SaveUserSession(context)
	fmt.Println("User logged out: " + session.Username)

	return context.Redirect(http.StatusFound, router.Routes.Login)
}

func About(w http.ResponseWriter, r *http.Request) {
	msg := "Hello from About, World!"
	w.Write([]byte(msg))
}
