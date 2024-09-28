package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/cmd/handlers"
)

func main() {
	data := &handlers.MyDB{}
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	data.MyData = db
	http.HandleFunc("/", data.HomePage)
	http.HandleFunc("/login", data.LoginPage)
	http.HandleFunc("/register", data.RegisterPage)
	http.ListenAndServe(":8080", nil)
}
