package main

import (
	"log"

	"ovo-server/internal/config"
	"ovo-server/internal/controller"
	"ovo-server/internal/database"
	customMiddleware "ovo-server/internal/middleware"
	"ovo-server/internal/model"
	"ovo-server/internal/router"
	"ovo-server/internal/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	log.Println("Initializing OVO Server... This build is for development purposes only.")
	// Initializing configuration, reading .env file and setting up environment variables
	config.Init()
	// Initializing database
	database.Init()
	// Migrating every time we start the server, this should be addressed to check versioning of database
	model.Init()
	// Session setup
	session.GenerateSessionHandler("TODO:TEMPORAL_COOKIE_SECRET_MUST_CHANGE", "ovo-session")
}

func main() {
	log.Println("Starting OVO Server...")
	router.InitRoutes()
	echoInstance := echo.New()
	// Static files route setup
	echoInstance.Static("/assets", "public")

	// Middleware setup
	// 		CORS setup
	echoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:1234"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// 		Request log setup
	echoInstance.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} - STATUS: ${status} - Method: ${method} - URI: ${uri}\n",
	}))

	// Route definition
	// 		Base Path
	echoBasePath := echoInstance.Group(router.GetBasePath())

	// 		Unauthenticated routes (Public routes)
	echoUnauthenticateGroup := echoBasePath.Group("")
	echoUnauthenticateGroup.Use(customMiddleware.IsNotAuthenticated)
	echoUnauthenticateGroup.GET(router.Routes.Login, controller.Login)
	echoUnauthenticateGroup.POST(router.Routes.Login, controller.LoginRequest)
	echoUnauthenticateGroup.GET(router.Routes.Register, controller.Register)
	echoUnauthenticateGroup.POST(router.Routes.Register, controller.RegisterRequest)

	//   	Authenticated routes (Private routes)
	// 			Visitor routes (unprivileged user)
	echoAuthenticatedGroup := echoBasePath.Group("")
	echoAuthenticatedGroup.Use(customMiddleware.IsAuthenticated)
	echoAuthenticatedGroup.GET(router.Routes.Logout, controller.Logout)
	echoAuthenticatedGroup.GET(router.Routes.Home, controller.Home)

	// 			Admin routes (admin only)
	echoAdminGroup := echoBasePath.Group("")
	echoAdminGroup.Use(customMiddleware.IsAdmin, customMiddleware.IsAuthenticated)
	echoAdminGroup.GET(router.AdminRoutes.Dashboard, controller.AdminDashboard)
	echoAdminGroup.GET(router.AdminRoutes.Libraries, controller.AdminLibraries)

	echoInstance.Start("localhost:8080")
}
