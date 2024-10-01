package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/cmd/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fs := http.FileServer(http.Dir("./cmd/templates/visualize"))
	http.Handle("/cmd/templates/visualize/", http.StripPrefix("/cmd/templates/visualize/", fs))
	data := &handlers.MyDB{}
	db, err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
	stm, err := db.Prepare("CREATE TABLE IF NOT EXISTS login (id INTEGER PRIMARY KEY AUTOINCREMENT, user TEXT NOT NULL, pass TEXT NOT NULL, uid TEXT)")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	stm.Exec()
	stm, err = db.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT ,user TEXT, post TEXT, comments TEXT, likes INTEGER DEFAULT 0, deslikes INTEGER DEFAULT 0)")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	stm.Exec()
	data.MyData = db
	http.HandleFunc("/", data.HomePage)
	http.HandleFunc("/login", data.LoginPage)
	http.HandleFunc("/register", data.RegisterPage)
	http.HandleFunc("/logout", data.Logout)
	http.HandleFunc("/post", data.PostsUser)
	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
