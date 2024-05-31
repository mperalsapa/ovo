package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"ovo-server/internal/config"
	"ovo-server/internal/controller"
	apiController "ovo-server/internal/controller/api"
	"ovo-server/internal/database"
	"ovo-server/internal/file"
	customMiddleware "ovo-server/internal/middleware"
	"ovo-server/internal/model"
	"ovo-server/internal/router"
	"ovo-server/internal/session"
	"ovo-server/internal/syncplay"
	"ovo-server/internal/tmdb"
	"ovo-server/internal/websocket"

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
	// Websocket setup
	websocket.Init()
	// Syncplay setup
	syncplay.Init()

	// Initialize TMDB API
	tmdb.Init()
}

//go:embed public
var staticAssets embed.FS

func main() {
	log.Println("Starting OVO Server...")
	echoInstance := echo.New()

	// Hide banner and port
	echoInstance.HideBanner = true
	echoInstance.HidePort = true

	// Static files route setup
	// echoInstance.Static(router.Routes.Assets, "public")
	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(file.GetFileSystem(useOS, staticAssets))
	echoInstance.GET(router.Routes.Assets+"/*", echo.WrapHandler(http.StripPrefix(router.Routes.Assets+"/", assetHandler)))

	// Middleware setup
	// 		CORS setup
	echoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:1234"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// 		Request log setup
	echoInstance.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} - IP: ${remote_ip} - STATUS: ${status} - Method: ${method} - URI: ${uri}\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
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
	auth.Use(customMiddleware.IsAuthenticated, customMiddleware.UpdateDeviceActivity, customMiddleware.ValidateCurrentSyncplayGroup)
	auth.GET(router.Routes.Logout, controller.Logout)
	auth.GET(router.Routes.Home, controller.Home)
	auth.GET(router.Routes.Library, controller.Library)
	auth.GET(router.Routes.Item, controller.Item)
	auth.GET(router.Routes.Person, controller.Person)
	auth.GET(router.Routes.Player, controller.Player)
	auth.GET(router.Routes.DownloadItem, controller.Download)
	auth.GET(router.Routes.Websocket, controller.WebsocketHandler)

	// 			API routes
	api := auth.Group("")
	api.Use(customMiddleware.IsAuthenticated, customMiddleware.UpdateDeviceActivity, customMiddleware.ValidateCurrentSyncplayGroup)
	api.POST(router.ApiRoutes.SyncplayGroups, apiController.CreateSyncGroup)
	api.GET(router.ApiRoutes.SyncplayGroups, apiController.GetSyncGroups)
	api.DELETE(router.ApiRoutes.SyncplayGroups, apiController.LeaveSyncGroup)
	api.PUT(router.ApiRoutes.SyncplayGroups, apiController.JoinSyncGroup)
	api.POST(router.ApiRoutes.ToggleFavoriteItem, apiController.ToggleFavoriteItem)
	api.POST(router.ApiRoutes.ToggleWatchedItem, apiController.ToggleWatchedItem)

	// 			Admin routes (admin only)
	admin := echoInstance.Group("")
	admin.Use(customMiddleware.IsAuthenticated, customMiddleware.IsAdmin, customMiddleware.UpdateDeviceActivity)
	admin.GET(router.AdminRoutes.Dashboard, controller.AdminDashboard)
	admin.GET(router.AdminRoutes.Libraries, controller.AdminLibraries)
	admin.GET(router.AdminRoutes.Library, controller.AdminLibrary)
	admin.POST(router.AdminRoutes.Library, controller.AdminStoreLibrary)
	admin.GET(router.AdminRoutes.Command, controller.AdminCommand)

	// DEBUG Print current echo routes
	// for _, route := range echoInstance.Routes() {
	// 	log.Println(route.Method, route.Path)
	// }

	log.Printf("Ready to serve requests. Started on :8080%s", router.BasePath)
	echoInstance.Start(":8080")
}
