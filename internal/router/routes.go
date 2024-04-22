package router

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

type route struct {
	Assets   string
	Api      string
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
	Command   string
}

type apiRoutes struct {
	Library   string
	Libraries string
}

var Routes route
var ApiRoutes apiRoutes
var AdminRoutes adminRoute

var BasePath = "/ovo"
var AdminBasePath = "/admin"

// func BuildRoute(path string) (string, error) {
// 	return url.JoinPath(BasePath, path)
// }

func BuildRoute(path string) string {
	route, err := url.JoinPath(BasePath, path)
	log.Println("Building route: ", route, " with path: ", path, " and basepath: ", BasePath)
	if err != nil {
		log.Println(err)
	}
	return route
}

func BuildApiRoute(path string) string {
	path, err := url.JoinPath(Routes.Api, path)
	if err != nil {
		log.Println(err)
	}
	return path
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

// generate route with id
func GenerateRouteWithId(route string, id uint) string {
	return strings.ReplaceAll(route, ":id", fmt.Sprintf("%d", id))
}

// generate route with string parameter
func GenerateRouteWithCommand(route string, param string) string {
	return strings.ReplaceAll(route, ":action", param)
}

func InitRoutes() bool {
	Routes.Assets = BuildRoute("/assets")
	Routes.Api = BuildRoute("/api")
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
	AdminRoutes.Settings = BuildAdminRoute("/settings")
	AdminRoutes.Command = BuildAdminRoute("/command/:action")

	// Api routes
	ApiRoutes.Library = BuildApiRoute("/library/:id")
	ApiRoutes.Libraries = BuildApiRoute("/libraries")

	return true
}
