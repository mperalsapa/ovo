package middleware

import (
	"log"
	"net/http"
	"ovo-server/internal/model"
	"ovo-server/internal/router"
	"ovo-server/internal/session"

	"github.com/labstack/echo/v4"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := session.IsAuth(c)
		if !auth {
			return c.Redirect(http.StatusFound, "/login")
		}

		UserExists(next)(c)

		return next(c)
	}
}

func IsNotAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := session.IsAuth(c)
		if auth {
			return c.Redirect(http.StatusFound, router.Routes.Home)
		}
		return next(c)
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := session.GetUserSession(c)
		if session.Username == "" {
			return c.Redirect(http.StatusFound, router.Routes.Login)
		}

		if session.Role != model.Admin {
			return c.Redirect(http.StatusFound, router.Routes.Home)
		}
		return next(c)
	}
}

func UserExists(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := session.GetUserSession(c)
		// Checking if user exists in the database
		exists := model.UserExists(user.Username)
		if !exists {
			log.Println("Invalid user session. User does not exist in the database.")
			user.Authenticated = false
			user.SaveUserSession(c)
			return c.Redirect(http.StatusFound, router.Routes.Login)
		}
		return next(c)
	}
}
