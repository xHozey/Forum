package handlers

import (
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/index.html")
	if err != nil {
		http.Error(w, "lkjdflkjds", 500)
		return
	}

	tmp.Execute(w, nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
}
