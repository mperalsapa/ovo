package main

import (
	"fmt"

	"ovo-server/internal/config"
	"ovo-server/internal/controller"
	"ovo-server/internal/database"
	customMiddleware "ovo-server/internal/middleware"
	model "ovo-server/internal/model"
	"ovo-server/internal/router"
	customSession "ovo-server/internal/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	fmt.Println("Init...")
	config.Init()
	database.Init()
	// Migrating every time we start the server, this should be addressed to check versioning of database
	model.Init()
}

func main() {
	fmt.Println("Hello, World!")
	router.InitRoutes()
	echoInstance := echo.New()

	echoInstance.Static("/assets", "public")
	echoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:1234"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	echoInstance.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} - STATUS: ${status} - Method: ${method} - URI: ${uri}\n",
	}))

	// staticFilesRoute := "../ovo-web/dist"
	// echoInst.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	// 	Root:   staticFilesRoute,
	// 	Index:  "index.html",
	// 	Browse: false,
	// 	HTML5:  true,
	// }))

	// echoInst.Use(customMiddleware.AuthMiddleware)
	customSession.GenerateSessionHandler("TODO:TEMPORAL_COOKIE_SECRET_MUST_CHANGE", "ovo-session")

	// Route definition
	echoInstance.GET(router.Routes.Login, controller.Login, customMiddleware.IsNotAuthenticated)
	echoInstance.POST(router.Routes.Login, controller.LoginRequest, customMiddleware.IsNotAuthenticated)
	echoInstance.GET(router.Routes.Register, controller.Register, customMiddleware.IsNotAuthenticated)
	echoInstance.POST(router.Routes.Register, controller.RegisterRequest, customMiddleware.IsNotAuthenticated)

	echoAuthenticatedGroup := echoInstance.Group("")
	echoAuthenticatedGroup.Use(customMiddleware.IsAuthenticated)
	echoAuthenticatedGroup.GET(router.Routes.Logout, controller.Logout)
	echoAuthenticatedGroup.GET(router.Routes.Home, controller.Home)
	echoAuthenticatedGroup.GET("/setpassword", controller.SetPassword)

	echoAdminGroup := echoInstance.Group("/admin")
	echoAdminGroup.Use(customMiddleware.IsAdmin)
	echoAdminGroup.GET("/", controller.AdminDashboard)
	// TODO IMPLEMENT ADMIN ROUTES

	echoInstance.Start("localhost:8080")
}
