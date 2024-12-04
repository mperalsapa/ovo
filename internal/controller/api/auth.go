package apiController

import (
	"log"
	"net/http"
	"ovo-server/internal/model"
	"ovo-server/internal/session"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mileusna/useragent"
)

type ApiResponse struct {
	Data  string `json:",omitempty"`
	Error string `json:",omitempty"`
}

func CheckAuth(context echo.Context) error {
	isAuth := session.IsAuth(context)

	return context.JSON(http.StatusOK, &ApiResponse{Data: strconv.FormatBool(isAuth)})
}

func Login(context echo.Context) error {
	isAuth := session.IsAuth(context)
	if isAuth {
		return context.JSON(http.StatusOK, &ApiResponse{Data: "Already logged in"})
	}

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
		userSession.SaveUserSession(context)
		return context.JSON(http.StatusUnauthorized, &ApiResponse{Error: "Invalid username or password."})
	}

	if !user.Enabled {
		userSession.Authenticated = false
		userSession.ErrorMsg = "This user is disabled"
		userSession.SaveUserSession(context)
		return context.JSON(http.StatusUnauthorized, &ApiResponse{Error: "This user is disabled"})
	}

	parsedUA := useragent.Parse(context.Request().UserAgent())
	device := model.CreateDevice(model.GetUserByUsername(userSession.Username).ID, parsedUA.Name)
	userSession.DeviceID = device.ID

	userSession.Authenticated = true
	userSession.ErrorMsg = ""
	userSession.Role = user.Role
	userSession.SaveUserSession(context)

	return context.JSON(http.StatusOK, &ApiResponse{Data: "Login Succesful"})
}

type RegisterRequest struct {
	Username             string `form:"username"`
	Password             string `form:"password"`
	PasswordVerification string `form:"password_verification"`
}

func Register(context echo.Context) error {
	userSession := session.GetUserSession(context)
	var reqUser RegisterRequest
	if err := context.Bind(&reqUser); err != nil {
		// return context.JSON(http.StatusBadRequest, err)
		return context.JSON(http.StatusBadRequest, &ApiResponse{Error: err.Error()})
	}

	userSession.Username = reqUser.Username
	if model.GetUserExists(reqUser.Username) {
		log.Println("User already exists: " + reqUser.Username)
		userSession.ErrorMsg = "User already exists"
		userSession.SaveUserSession(context)
		// return context.Redirect(http.StatusFound, router.Routes.Register)
		return context.JSON(http.StatusBadRequest, &ApiResponse{Error: "User already exists"})
	}

	if reqUser.Password != reqUser.PasswordVerification {
		log.Println("Password and password verification do not match")
		userSession.ErrorMsg = "Password and password verification do not match"
		userSession.SaveUserSession(context)
		// return context.Redirect(http.StatusFound, router.Routes.Register)
		return context.JSON(http.StatusBadRequest, &ApiResponse{Error: "Password and password verification do not match"})
	}

	if reqUser.Password == "" {
		return context.JSON(http.StatusBadRequest, &ApiResponse{Error: "Password cannot be empty"})
	}

	user := model.NewUser(reqUser.Username, reqUser.Password)

	// Disables new users if there are already users in the system
	// This was added to prevent users from registering without the admin
	// consent. Commented temporarlly to allow users to register and
	// in a future we will implement the admin consent in admin dashboard.
	userCount := model.UserCount()
	if userCount == 0 {
		user.Role = model.Admin
		user.Enabled = true
	} else {
		user.Enabled = true
	}

	user.Save()

	// return context.Redirect(http.StatusFound, router.Routes.Login)
	return context.JSON(http.StatusOK, &ApiResponse{Data: "true"})
}

func Logout(context echo.Context) error {
	session := session.GetUserSession(context)
	session.Authenticated = false
	session.SaveUserSession(context)
	log.Println("User logged out: " + session.Username)
	return context.JSON(http.StatusOK, &ApiResponse{Data: "true"})
}
