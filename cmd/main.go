package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/cmd/handlers"
)

func main() {
	data := &handlers.MyDB{}
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	stm, err := db.Prepare("CREATE TABLE IF NOT EXISTS login (id INTEGER PRIMARY KEY AUTOINCREMENT, user TEXT, pass TEXT)")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	stm.Exec()
	data.MyData = db
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/login", data.LoginPage)
	http.HandleFunc("/register", data.RegisterPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
