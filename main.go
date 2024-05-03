package main

import (
	"fmt"
	"log"

	"ovo-server/internal/config"
	"ovo-server/internal/controller"
	"ovo-server/internal/database"
	customMiddleware "ovo-server/internal/middleware"
	"ovo-server/internal/model"
	"ovo-server/internal/router"
	"ovo-server/internal/session"
	"ovo-server/internal/tmdb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	// Display banner

	banner := `
	  ____ _   ______ 
	 / __ \ | / / __ \
	/ /_/ / |/ / /_/ /
	\____/|___/\____/ 
	`
	fmt.Println(banner)
	log.Println("Initializing OVO Server... This build is for development purposes only.")
	// Initializing configuration, reading .env file and setting up environment variables
	config.Init()
	// Initializing database
	database.Init()
	// Migrating every time we start the server, this should be addressed to check versioning of database
	model.Init()
	// Session setup
	session.GenerateSessionHandler("TODO:TEMPORAL_COOKIE_SECRET_MUST_CHANGE", "ovo-session")
	// Router setup
	router.Init()

	// Initialize TMDB API
	tmdb.Init()
}

func main() {
	log.Println("Starting OVO Server...")
	echoInstance := echo.New()

	// Hide banner and port
	echoInstance.HideBanner = true
	echoInstance.HidePort = true

	// Static files route setup
	echoInstance.Static(router.Routes.Assets, "public")

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
	// 		Unauthenticated routes (Public routes)
	unauth := echoInstance.Group("")
	unauth.Use(customMiddleware.IsNotAuthenticated)
	unauth.GET(router.Routes.Login, controller.Login)
	unauth.POST(router.Routes.Login, controller.LoginRequest)
	unauth.GET(router.Routes.Register, controller.Register)
	unauth.POST(router.Routes.Register, controller.RegisterRequest)

	//   	Authenticated routes (Private routes)
	// 			Visitor routes (unprivileged user)
	auth := echoInstance.Group("")
	auth.Use(customMiddleware.IsAuthenticated)
	auth.GET(router.Routes.Logout, controller.Logout)
	auth.GET(router.Routes.Home, controller.Home)
	auth.GET(router.Routes.Library, controller.Library)
	auth.GET(router.Routes.Item, controller.Item)

	// 			API routes
	api := auth.Group("")
	api.GET(router.ApiRoutes.Libraries, controller.APIGetLibraries)
	api.GET(router.ApiRoutes.Library, controller.APIGetLibrary)
	api.POST(router.ApiRoutes.Library, controller.APIAddLibrary)
	api.DELETE(router.ApiRoutes.Library, controller.APIDeleteLibrary)

	// 			Admin routes (admin only)
	admin := echoInstance.Group("")
	admin.Use(customMiddleware.IsAdmin, customMiddleware.IsAuthenticated)
	admin.GET(router.AdminRoutes.Dashboard, controller.AdminDashboard)
	admin.GET(router.AdminRoutes.Libraries, controller.AdminLibraries)
	admin.GET(router.AdminRoutes.Library, controller.AdminLibrary)
	admin.POST(router.AdminRoutes.Library, controller.AdminStoreLibrary)
	admin.GET(router.AdminRoutes.Command, controller.AdminCommand)

	// DEBUG Print current echo routes
	// for _, route := range echoInstance.Routes() {
	// 	log.Println(route.Method, route.Path)
	// }

	log.Println("Started on http://localhost:8080. Ready to serve requests.")
	echoInstance.Start("localhost:8080")
}
