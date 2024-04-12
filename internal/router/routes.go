package router

type route struct {
	Login    string
	Logout   string
	Register string
	Home     string
	Library  string
	Profile  string
}

var Routes route
var BasePath = "/"

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
}
