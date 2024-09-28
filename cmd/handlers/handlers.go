package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func (d *MyDB) HomePage(w http.ResponseWriter, r *http.Request) {

}

func (d *MyDB) RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/register.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")
	statement, err := d.MyData.Prepare("INSERT INTO login (mail, password) VALUES (?,?)")
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")
	statement, err := d.MyData.Query("SELECT (mail, password) FROM person")
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
