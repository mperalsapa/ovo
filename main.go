package main

import (
	"fmt"

	"ovo-server/internal/controller"
	"ovo-server/internal/database"
	model "ovo-server/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func middleLogger1(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Middleware Logger")
		return next(c)
	}
}

func middleLogger2(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Middleware Logger 2")
		return next(c)
	}
}

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
	echoInst.Static("/", staticFilesRoute)

	echoInst.POST("/login", controller.Login)

	apiGroup := echoInst.Group("/api")
	apiGroup.GET("/home", controller.Home, middleLogger1, middleLogger2)
	apiGroup.GET("/setpassword", controller.SetPassword)
	apiGroup.GET("/register", controller.Register)

	echoInst.Start("localhost:8080")
}
