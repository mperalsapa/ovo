package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"ovo-server/internal/config"
	"strings"
)

type route struct {
	Assets          string
	Api             string
	Login           string
	Logout          string
	Register        string
	Profile         string
	Home            string
	Library         string
	FavoriteLibrary string
	Item            string
	Person          string
	Player          string
	DownloadItem    string
	Websocket       string
	RoutesJSON      string
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
	Library            string
	Libraries          string
	SyncplayGroups     string // CRUD = POST, GET, DELETE, PUT
	ToggleFavoriteItem string // POST
	ToggleWatchedItem  string // POST
}

var Routes route
var ApiRoutes apiRoutes
var AdminRoutes adminRoute

// Default basepath for app is "/". This can be changed in the .env file.
var BasePath = "/"
var AdminBasePath = "/admin"

// func BuildRoute(path string) (string, error) {
// 	return url.JoinPath(BasePath, path)
// }

func BuildRoute(path string) string {
	route, err := url.JoinPath(BasePath, path)
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

func GeneratePlayerRoute(itemID uint) string {
	return fmt.Sprintf("%s?item=%d", Routes.Player, itemID)
}

func GenerateDownloadItemRoute(itemID uint) string {
	return fmt.Sprintf("%s?item=%d", Routes.DownloadItem, itemID)
}

func GetRoutesJSON() string {
	routesJSON, err := json.Marshal(map[string]interface{}{
		"Routes":    Routes,
		"ApiRoutes": ApiRoutes,
	})

	if err != nil {
		log.Println("Error marshalling routes to JSON:", err)
		return ""
	}

	return string(routesJSON)
}

func SaveRoutesJSON() {

	routesJSON := GetRoutesJSON()

	if routesJSON == "" {
		return
	}

	err := os.WriteFile("public/routes.json", []byte(routesJSON), 0644)

	if err != nil {
		log.Println("Error saving routes to file:", err)
		return
	}
}

func Init() bool {

	if config.Variables.Basedir != "" && config.Variables.Basedir != "/" {
		BasePath = config.Variables.Basedir
	}

	Routes.Assets = BuildRoute("/assets")
	Routes.Api = BuildRoute("/api")
	Routes.Login = BuildRoute("/login")
	Routes.Logout = BuildRoute("/logout")
	Routes.Register = BuildRoute("/register")
	Routes.Profile = BuildRoute("/profile")
	Routes.Home = BuildRoute("/")
	Routes.Library = BuildRoute("/library/:id")
	Routes.FavoriteLibrary = BuildRoute("/favorite-library")
	Routes.Item = BuildRoute("/item/:id")
	Routes.Person = BuildRoute("/person/:id")
	Routes.Player = BuildRoute("/player")
	Routes.DownloadItem = BuildRoute("/download")
	Routes.Websocket = BuildRoute("/ws")
	Routes.RoutesJSON = BuildRoute("/assets/routes.json")

	// Admin routes
	AdminRoutes.Dashboard = BuildAdminRoute("")
	AdminRoutes.Libraries = BuildAdminRoute("/libraries")
	AdminRoutes.Library = BuildAdminRoute("/library/:id")
	AdminRoutes.Users = BuildAdminRoute("/users")
	AdminRoutes.User = BuildAdminRoute("/user/:id")
	AdminRoutes.Settings = BuildAdminRoute("/settings")
	AdminRoutes.Command = BuildAdminRoute("/command/:action")

	// Api routes
	// ApiRoutes.Library = BuildApiRoute("/library/:id")
	// ApiRoutes.Libraries = BuildApiRoute("/libraries")
	ApiRoutes.SyncplayGroups = BuildApiRoute("/syncplay/groups")
	ApiRoutes.ToggleFavoriteItem = BuildApiRoute("/toggle-favorite-item")
	ApiRoutes.ToggleWatchedItem = BuildApiRoute("/toggle-watched-item")

	SaveRoutesJSON()

	return true
}
