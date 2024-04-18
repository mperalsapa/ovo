package router

import (
	"fmt"
	"log"
	"net/url"
)

type route struct {
	Assets   string
	Login    string
	Logout   string
	Register string
	Home     string
	Library  string
	Profile  string
}

type adminRoute struct {
	Dashboard string
	Libraries string
	Library   string
	Users     string
	User      string
	Settings  string
}

var Routes route
var AdminRoutes adminRoute
var BasePath = "/ovo"
var AdminBasePath = "/admin"

// func BuildRoute(path string) (string, error) {
// 	return url.JoinPath(BasePath, path)
// }

func BuildRoute(path string) string {
	route, err := url.JoinPath(BasePath, path)
	fmt.Println("Building route: ", route, " with path: ", path, " and basepath: ", BasePath)
	if err != nil {
		log.Println(err)
	}
	return route
}

func BuildAdminRoute(path string) string {
	path, err := url.JoinPath(BasePath, AdminBasePath, path)
	if err != nil {
		log.Println(err)
	}
	return path
}

func GetBasePath() string {
	if BasePath == "/" {
		return ""
	}

	return BasePath
}

func InitRoutes() {
	Routes.Assets = BuildRoute("/assets")
	Routes.Login = BuildRoute("/login")
	Routes.Logout = BuildRoute("/logout")
	Routes.Register = BuildRoute("/register")
	Routes.Home = BuildRoute("/")
	Routes.Library = BuildRoute("/library/:id")
	Routes.Profile = BuildRoute("/profile")

	// Admin routes
	AdminRoutes.Dashboard = BuildAdminRoute("")
	AdminRoutes.Libraries = BuildAdminRoute("/libraries")
	AdminRoutes.Library = BuildAdminRoute("/library/:id")
	AdminRoutes.Users = BuildAdminRoute("/users")
	AdminRoutes.User = BuildAdminRoute("/user/:id")
}
