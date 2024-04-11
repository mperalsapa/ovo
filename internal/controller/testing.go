package controller

import (
	"ovo-server/internal/model"

	"github.com/labstack/echo/v4"
)

func SetPassword(c echo.Context) error {
	pwd := c.QueryParam("pwd")
	user := model.User{}
	user.SetPassword(pwd)
	return c.JSON(200, user.Password)
}

func Register(c echo.Context) error {
	username := c.QueryParam("username")
	pwd := c.QueryParam("pwd")
	role := c.QueryParam("role")
	user := model.CreateUser(username, pwd, model.Role(role))
	user.Save()
	return c.JSON(200, user)
}

func LoginTest(c echo.Context) error {
	username := c.QueryParam("username")
	user := model.GetUserByUsername(username)
	return c.JSON(200, user)
}
