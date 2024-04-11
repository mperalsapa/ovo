package controller

import (
	"fmt"
	"net/http"
	"ovo-server/internal/model"

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

	return e.JSON(http.StatusOK, user)
}

func About(w http.ResponseWriter, r *http.Request) {
	msg := "Hello from About, World!"
	w.Write([]byte(msg))
}
