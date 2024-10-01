package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func (d *MyDB) authorize(r *http.Request) (bool, string) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return false, ""
	}

	var uid string
	var username string
	err = d.MyData.QueryRow("SELECT uid, user FROM login WHERE uid = ?", c.Value).Scan(&uid, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ""
		}
		fmt.Println("Database error:", err)
		return false, ""
	}

	return c.Value == uid, username
}

func (d *MyDB) insertPost(post string, user string) error {
	stm, err := d.MyData.Prepare("INSERT INTO posts (post, user) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(post, user)
	if err != nil {
		return err
	}

	return nil
}

func (d *MyDB) getPosts(user string) ([]Post, error) {
	rows, err := d.MyData.Query("SELECT user, id, likes, deslikes, post FROM posts WHERE user = ?", user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.User, &post.Id, &post.Like, &post.Deslike, &post.Post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (d *MyDB) HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./cmd/templates/index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	authorized, username := d.authorize(r)
	if !authorized {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		post := r.FormValue("textarea")
		if post != "" {
			err := d.insertPost(post, username)
			if err != nil {
				fmt.Print(err)
				http.Error(w, "Failed to insert post", http.StatusInternalServerError)
				return
			}
		}
	}

	posts, err := d.getPosts(username)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	data := struct {
		Authorized bool
		Username   string
		Posts      []Post
	}{
		Authorized: authorized,
		Username:   username,
		Posts:      posts,
	}

	tmp.Execute(w, data)
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
	if err != nil && err != http.ErrNoCookie {
		log.Fatal(err)
	}
	_, err = d.MyData.Exec("UPDATE login SET uid = '' WHERE uid = ?", c.Value)
	if err != nil {
		log.Printf("Database error during logout: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func generateUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func setSessionCookie(w http.ResponseWriter, uid string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  uid,
		Path:   "/",
		MaxAge: 3600,
	})
}
