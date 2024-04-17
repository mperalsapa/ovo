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

func InitRoutes() {
	Routes.Login = BuildRoute("/login")
	Routes.Logout = BuildRoute("/logout")
	Routes.Register = BuildRoute("/register")
	Routes.Home = BuildRoute("/")
	Routes.Library = BuildRoute("/library/:id")
	Routes.Profile = BuildRoute("/profile")

	// Admin routes
	AdminRoutes.Dashboard = BuildRoute(AdminBasePath)
	AdminRoutes.Libraries = BuildRoute(AdminBasePath + "/libraries")
	AdminRoutes.Library = BuildRoute(AdminBasePath + "/library/:id")
	AdminRoutes.Users = BuildRoute(AdminBasePath + "/users")
	AdminRoutes.User = BuildRoute(AdminBasePath + "/user/:id")
}
