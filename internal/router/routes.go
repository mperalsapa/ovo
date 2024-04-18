package router

type route struct {
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
var BasePath = "/"
var AdminBasePath = "/admin"

func BuildRoute(path string) string {
	// if basepath is root, return the path
	if BasePath == "/" {
		return path
	}

	// check that basepath begins with a slash
	if BasePath[0] != '/' || BasePath == "" {
		BasePath = "/" + BasePath
	}

	// return the basepath concatenated with the path
	return BasePath + path
}

func BuildAdminRoute(path string) string {
	// ensuring path contains / at the start
	if path[0] != '/' {
		path = "/" + path
	}

	// check that basepath begins with a slash
	if AdminBasePath[0] != '/' || AdminBasePath == "" {
		AdminBasePath = "/" + AdminBasePath
	}

	// if basepath is root, return the path with hardcoded /admin
	// this is to prevent double handlers on the same route
	if AdminBasePath == "/" {
		return "/admin" + path
	}

	// return the basepath concatenated with the path
	return AdminBasePath + path
}

func GetBasePath() string {
	if BasePath == "/" {
		return ""
	}

	return BasePath
}

func InitRoutes() {
	Routes.Login = "/login"
	Routes.Logout = "/logout"
	Routes.Register = "/register"
	Routes.Home = "/"
	Routes.Library = "/library/:id"
	Routes.Profile = "/profile"

	// Admin routes
	AdminRoutes.Dashboard = AdminBasePath
	AdminRoutes.Libraries = AdminBasePath + "/libraries"
	AdminRoutes.Library = AdminBasePath + "/library/:id"
	AdminRoutes.Users = AdminBasePath + "/users"
	AdminRoutes.User = AdminBasePath + "/user/:id"
}
