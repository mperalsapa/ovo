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
		if !auth || !UserExist(c) {
			return c.Redirect(http.StatusFound, router.Routes.Login)
		}

		return next(c)
	}
}

func IsNotAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := session.IsAuth(c)
		if auth {
			log.Println("User is already authenticated, redirecting to home: ", router.Routes.Home)
			return c.Redirect(http.StatusFound, router.Routes.Home)
		}
		return next(c)
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := session.GetUserSession(c)
		if session.Username == "" {
			log.Println("No session found, redirecting to login: ", router.Routes.Login)
			return c.Redirect(http.StatusFound, router.Routes.Login)
		}

		if session.Role != model.Admin {
			log.Println("User is not an admin, redirecting to home: ", router.Routes.Home)
			return c.Redirect(http.StatusFound, router.Routes.Home)
		}
		return next(c)
	}
}

func UserExist(echo echo.Context) bool {
	userSession := session.GetUserSession(echo)
	user := model.GetUserByUsername(userSession.Username)
	// exists := model.GetUserExists(userSession.Username)
	if user.Username == "" || !user.Enabled {
		userSession.Authenticated = false
		userSession.SaveUserSession(echo)
		return false
	}

	return true
}
