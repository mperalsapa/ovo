package controller

import (
	"fmt"
	"log"
	"net/http"
	"ovo-server/internal/model"
	"ovo-server/internal/router"
	"ovo-server/internal/session"
	"ovo-server/internal/template/page"

	"github.com/labstack/echo/v4"
)

func Login(context echo.Context) error {
	userSession := session.GetUserSession(context)
	pageData := page.LoginPageData{
		Username: userSession.Username,
		AlertMsg: userSession.PopErrorMessage(context),
	}

	component := page.LoginPage(pageData)
	return RenderView(context, http.StatusOK, component)
}

func LoginRequest(context echo.Context) error {
	var reqUser model.User
	if err := context.Bind(&reqUser); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	userSession := session.GetUserSession(context)
	userSession.Username = reqUser.Username
	log.Println("Login request for user " + reqUser.Username)
	user := model.GetUserByUsername(reqUser.Username)

	if valid := user.CheckPassword(reqUser.Password); !valid {
		userSession.Authenticated = false
		userSession.ErrorMsg = "Invalid username or password"
		userSession.SaveUserSession(context)
		return context.Redirect(http.StatusFound, router.Routes.Login)
	}

	if !user.Enabled {
		userSession.Authenticated = false
		userSession.ErrorMsg = "This user is disabled"
		userSession.SaveUserSession(context)
		return context.Redirect(http.StatusFound, router.Routes.Login)
	}

	userSession.Authenticated = true
	userSession.ErrorMsg = ""
	userSession.Role = user.Role
	userSession.SaveUserSession(context)

	return context.Redirect(http.StatusFound, router.Routes.Home)
}

func Register(context echo.Context) error {
	userSession := session.GetUserSession(context)
	log.Println("Register page - Current error msg: " + userSession.ErrorMsg)
	pageData := page.RegisterPageData{
		Username: userSession.Username,
		AlertMsg: userSession.PopErrorMessage(context),
	}
	fmt.Println("Register page - Register error msg: " + pageData.AlertMsg)
	component := page.RegisterPage(pageData)
	return RenderView(context, http.StatusOK, component)
}

func RegisterRequest(context echo.Context) error {
	userSession := session.GetUserSession(context)

	type RegisterRequest struct {
		Username             string `form:"username"`
		Password             string `form:"password"`
		PasswordVerification string `form:"password_verification"`
	}

	var reqUser RegisterRequest
	if err := context.Bind(&reqUser); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	userSession.Username = reqUser.Username
	if model.GetUserExists(reqUser.Username) {
		log.Println("User already exists: " + reqUser.Username)
		userSession.ErrorMsg = "User already exists"
		userSession.SaveUserSession(context)
		return context.Redirect(http.StatusFound, router.Routes.Register)
	}

	if reqUser.Password != reqUser.PasswordVerification {
		log.Println("Password and password verification do not match")
		userSession.ErrorMsg = "Password and password verification do not match"
		userSession.SaveUserSession(context)
		return context.Redirect(http.StatusFound, router.Routes.Register)
	}

	user := model.NewUser(reqUser.Username, reqUser.Password)

	userCount := model.UserCount()
	if userCount == 0 {
		user.Role = model.Admin
		user.Enabled = true
	}

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
