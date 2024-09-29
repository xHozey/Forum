package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func (d *MyDB) HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/index.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound) // Redirect to login if cookie not found
		return
	}

	qry, err := d.MyData.Query("SELECT uid FROM login WHERE uid = ?", c.Value)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer qry.Close()

	if qry.Next() {
		fmt.Println("youre good")
	} else {
		fmt.Println("session invalid")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tmp.Execute(w, nil)
}

func (d *MyDB) RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/register.htm")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("user")
	password := r.FormValue("password")
	if name != "" && password != "" {
		qry, err := d.MyData.Query("SELECT user FROM login WHERE user = ?", name)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer qry.Close()

		for qry.Next() {
			var user string
			err = qry.Scan(&user)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if user == name {
				http.Error(w, "user already exists", http.StatusConflict)
				return
			}
		}

		statement, err := d.MyData.Prepare("INSERT INTO login (user, pass, uid) VALUES (?, ?, ?)")
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer statement.Close()

		_, err = statement.Exec(name, password, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)

	}

	tmp.Execute(w, nil)
}

func (d *MyDB) LoginPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/login.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("user")
	password := r.FormValue("password")
	statement, err := d.MyData.Query("SELECT user, pass FROM login")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	for statement.Next() {
		var user, pass string
		statement.Scan(&user, &pass)
		if user == name && password == pass {
			u := uuid.Must(uuid.NewV4())
			uid := u.String()
			stm, err := d.MyData.Prepare("UPDATE login SET uid = ? WHERE user = ?")
			if err != nil {
				log.Fatal(err)
			}
			stm.Exec(uid, user)
			http.SetCookie(w, &http.Cookie{
				Name:   "session_token",
				Value:  uid,
				Path:   "/",
				MaxAge: 60,
			})

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
	}

	tmp.Execute(w, nil)
}
