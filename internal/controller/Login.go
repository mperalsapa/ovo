package controller

import (
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case "GET":
		LoginPage(w, r)
	case "POST":
		LoginForm(w, r)
	}

}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	msg := "Hello from LoginPage, World!"
	w.Write([]byte(msg))
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	msg := "Hello from LoginForm, World!"
	w.Write([]byte(msg))
}

func About(w http.ResponseWriter, r *http.Request) {
	msg := "Hello from About, World!"
	w.Write([]byte(msg))
}
