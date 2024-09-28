package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/index.html")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	tmp.Execute(w, nil)
}

func (d *MyDB) RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("/home/hamza/Desktop/Forum/cmd/templates/register.htm")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")
	statement, err := d.MyData.Prepare("INSERT INTO login (user, pass) VALUES (?,?)")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	statement.Exec(name, password)
	tmp.Execute(w, nil)
}

func (d *MyDB) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/index.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")
	statement, err := d.MyData.Query("SELECT (user, pass) FROM login")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for statement.Next() {
		var user, pass string
		statement.Scan(&user, &pass)
		if user == name && password == pass {
			http.Redirect(w, r, "/", 202)
		}
	}
	tmp.Execute(w, nil)
}
