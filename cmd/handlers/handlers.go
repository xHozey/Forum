package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func (d *MyDB) authorize(r *http.Request) bool {
	c, err := r.Cookie("session_token")
	if err != nil {
		return false
	}
	
	var uid string
	d.MyData.QueryRow("SELECT uid FROM login WHERE uid = ?", c.Value).Scan(&uid)
	return c.Value == uid
}

func (d *MyDB) HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	authorized := d.authorize(r)

	tmp.Execute(w, authorized)
}

func (d *MyDB) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		template.Must(template.ParseFiles("./cmd/templates/register.htm")).Execute(w, nil)
		return
	}

	name := r.FormValue("user")
	password := r.FormValue("password")

	if name == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	var exists bool
	err := d.MyData.QueryRow("SELECT 1 FROM login WHERE user = ?", name).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	uid := generateUID()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = d.MyData.Exec("INSERT INTO login (user, pass, uid) VALUES (?, ?, ?)", name, hashedPass, uid)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, uid)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (d *MyDB) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		template.Must(template.ParseFiles("./cmd/templates/login.html")).Execute(w, nil)
		return
	}

	name := r.FormValue("user")
	password := r.FormValue("password")

	var (
		storedPassword string
		uid            string
	)
	err := d.MyData.QueryRow("SELECT pass, uid FROM login WHERE user = ?", name).Scan(&storedPassword, &uid)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	newUID := generateUID()
	_, err = d.MyData.Exec("UPDATE login SET uid = ? WHERE user = ?", newUID, name)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, newUID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (d *MyDB) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		log.Fatal(err)
	}
	_, err = d.MyData.Exec("UPDATE login SET uid = '' WHERE uid = ?", c.Value)
	if err != nil {
		log.Printf("Database error during logout: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func generateUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func setSessionCookie(w http.ResponseWriter, uid string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    uid,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(3600 * time.Second),
	})
}
