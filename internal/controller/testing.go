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
