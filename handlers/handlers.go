package forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func (d *MyDB) HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	authorized, username := d.authorize(r)
	data := User{}
	data.Auth = authorized
	data.Username = username
	posts, err := d.getPosts()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
	data.Posts = posts
	like := r.FormValue("like")
	deslike := r.FormValue("deslike")
	var likeID, deslikeID int
	if like != "" {
		likeID, err = strconv.Atoi(like)
		if err != nil {
			fmt.Println(err)
		}
		d.insertLike(likeID, username)
	}
	if deslike != "" {
		deslikeID, err = strconv.Atoi(deslike)
		if err != nil {
			fmt.Println(err)
		}
		d.insertDeslike(deslikeID, username)
	}
	// fill our data variable with all posts
	for i := range posts {
		comments := d.getComment(data.Posts[i].Id)
		data.Posts[i].Comment = append(data.Posts[i].Comment, comments...)
		data.Posts[i].Like, data.Posts[i].Deslike = d.getLikes(data.Posts[i].Id)
	}
	if authorized {
		uid := refrechCooki(w)
		_, err = d.MyData.Exec("UPDATE login SET uid = ? WHERE user = ?", uid, username)
		if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	err = tmp.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *MyDB) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		template.Must(template.ParseFiles("./templates/register.htm")).Execute(w, nil)
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
		template.Must(template.ParseFiles("./templates/login.html")).Execute(w, nil)
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
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	newUID := generateUID()
	_, err = d.MyData.Exec("UPDATE login SET uid = ? WHERE user = ?", newUID, name)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, newUID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (d *MyDB) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil && err != http.ErrNoCookie {
		fmt.Println(err)
	}
	_, err = d.MyData.Exec("UPDATE login SET uid = '' WHERE uid = ?", c.Value)
	if err != nil {
		fmt.Println(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (d *MyDB) PostsUser(w http.ResponseWriter, r *http.Request) {
	post := r.FormValue("textarea")
	comment := r.FormValue("comment")
	postID := r.FormValue("postID")
	auth, username := d.authorize(r)

	if post != "" && auth {
		err := d.insertPost(post, username)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}
	}

	if comment != "" && postID != "" {
		id, err := strconv.Atoi(postID)
		if err != nil {
			fmt.Println(err)
		}
		err = d.insertComment(comment, id)
		if err != nil {
			fmt.Println(err)
		}

	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
