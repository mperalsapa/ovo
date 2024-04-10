package main

import (
	"fmt"

	"ovo-server/internal/controller"
	"ovo-server/internal/database"

	"github.com/labstack/echo/v4"
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
}

func main() {
	fmt.Println("Hello, World!")
	echoInst := echo.New()
	echoInst.GET("/", controller.Home, middleLogger1, middleLogger2)
	echoInst.GET("/setpassword", controller.SetPassword)

	echoInst.Start("localhost:8080")
}
