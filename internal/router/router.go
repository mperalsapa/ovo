package router

import (
	"fmt"
	"net/http"
	"ovo-server/internal/controller"
)

type Route struct {
	Name    string
	Path    string
	Handler http.HandlerFunc
}

func NewRoute(name string, path string, handler http.HandlerFunc) *Route {
	return &Route{
		Name:    name,
		Path:    path,
		Handler: handler,
	}
}

type Router struct {
	Routes map[string]*Route
}

func NewRouter() *Router {
	return &Router{
		Routes: make(map[string]*Route),
	}
}

func (r *Router) AddRoute(route *Route) {
	r.Routes[route.Name] = route
}

func (r *Router) GetRoute(name string) *Route {

	return r.Routes[name]
}

func (r *Router) RegisterRoutes() {
	for _, route := range r.Routes {
		// http.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
		// 	msg := fmt.Sprintf("Hello from %s, World!", route.Name)
		// 	w.Write([]byte(msg))
		// })
		// http.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {})
		http.HandleFunc(route.Path, route.Handler)
	}
}

var routes = NewRouter()
var PublicRoutes = NewRouter()

func init() {
	PublicRoutes.AddRoute(NewRoute("Login", "/login", controller.Login))
	// PublicRoutes.AddRoute(NewRoute("LoginForm", "POST", "/login", controller.LoginForm))
	PublicRoutes.AddRoute(NewRoute("About", "/about", controller.About))

	// routes.AddRoute(NewRoute("Hello", "GET", "/"))

	PublicRoutes.RegisterRoutes()
	routes.RegisterRoutes()
}

func Start() {
	fmt.Println("Starting web server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting web server: ", err)
	}
}
