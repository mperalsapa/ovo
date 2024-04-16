package middleware

import (
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
