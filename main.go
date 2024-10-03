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
	// Serve visualize to html
	fs := http.FileServer(http.Dir("./templates/visualize"))
	http.Handle("/templates/visualize/", http.StripPrefix("/templates/visualize/", fs))
	// struct that hold our data base to modify in our handlers
	data := &forum.MyDB{}
	db, err := sql.Open("sqlite3", "./db/Data.db")
	if err != nil {
		log.Fatal(err)
	}
	// function to create tables in Data.db if theyre not exists
	database.PrepareDataBase(db)
	data.MyData = db
	http.HandleFunc("/", data.HomePage)
	http.HandleFunc("/login", data.LoginPage)
	http.HandleFunc("/register", data.RegisterPage)
	http.HandleFunc("/logout", data.Logout)
	// to handle when user post
	http.HandleFunc("/post", data.PostsUser)
	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
