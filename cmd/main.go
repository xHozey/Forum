package main

import (
	"net/http"

	"forum/cmd/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/register", handlers.RegisterPage)
	http.ListenAndServe(":8080", nil)
}
