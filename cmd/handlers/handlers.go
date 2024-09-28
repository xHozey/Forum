package handlers

import (
	"fmt"
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
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodPost {
		fmt.Println("waaa adrari")
	}
	err := r.ParseForm() 
	if err != nil {
		fmt.Println(err)
		return
	}
	
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Println(name, email, password)
}
