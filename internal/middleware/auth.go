package middleware

import (
	"fmt"
	"net/http"
	"ovo-server/internal/model"
	localsession "ovo-server/internal/session"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Checking authentication v2 ...")
		r := c.Request()
		auth := localsession.IsAuth(r)
		if !auth {
			return c.Redirect(http.StatusFound, "/login")
		}

		return next(c)
	}
}

func IsNotAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Checking authentication v2 ...")
		r := c.Request()
		auth := localsession.IsAuth(r)
		if auth {
			return c.JSON(http.StatusOK, "Already authenticated")
		}

		return next(c)
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Checking admin")
		session, _ := session.Get("session", c)
		if session.Values["user"] == nil {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		user := session.Values["user"].(model.User)
		if user.Role != model.Admin {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}
