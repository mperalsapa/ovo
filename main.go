package main

import (
	"fmt"

	"ovo-server/internal/controller"
	"ovo-server/internal/database"
	customMiddleware "ovo-server/internal/middleware"
	model "ovo-server/internal/model"
	customSession "ovo-server/internal/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	fmt.Println("Init...")
	database.Init()
	// Migrating every time we start the server, this should be addressed to check versioning of database
	model.Init()
}

func main() {
	fmt.Println("Hello, World!")
	echoInst := echo.New()

	echoInst.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:1234"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	echoInst.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	staticFilesRoute := "../ovo-web/dist"
	echoInst.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   staticFilesRoute,
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))

	// echoInst.Use(customMiddleware.AuthMiddleware)
	customSession.GenerateSessionHandler("TODO:TEMPORAL_COOKIE_SECRET_MUST_CHANGE", "ovo-session")

	// Route definition
	echoInst.POST("/login", controller.Login, customMiddleware.IsNotAuthenticated)
	echoInst.POST("/register", controller.Register, customMiddleware.IsNotAuthenticated)
	apiGroup := echoInst.Group("/api")
	apiGroup.Use(customMiddleware.IsAuthenticated)
	apiGroup.GET("/logout", controller.Logout)
	apiGroup.GET("/home", controller.Home)
	apiGroup.GET("/setpassword", controller.SetPassword)
	apiGroup.GET("/register", controller.Register)

	echoInst.Start("localhost:8080")
}
