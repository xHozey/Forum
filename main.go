package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	database "forum/db"
	forum "forum/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fs := http.FileServer(http.Dir("./templates/visualize"))
	http.Handle("/templates/visualize/", http.StripPrefix("/templates/visualize/", fs))
	data := &forum.MyDB{}
	db, err := sql.Open("sqlite3", "./db/Data.db")
	if err != nil {
		log.Fatal(err)
	}
	database.PrepareDataBase(db)
	data.MyData = db
	http.HandleFunc("/", data.HomePage)
	http.HandleFunc("/login", data.LoginPage)
	http.HandleFunc("/register", data.RegisterPage)
	http.HandleFunc("/logout", data.Logout)
	http.HandleFunc("/post", data.PostsUser)
	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
