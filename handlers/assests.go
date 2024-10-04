package forum

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

// check if the user is logged in if he's logged we return he's username with true boolean else we return empty string with false
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

// generate unique random token for users
func generateUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

// give user a cookie
func setSessionCookie(w http.ResponseWriter, uid string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  uid,
		Path:   "/",
		MaxAge: 3600,
	})
}

func deleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

// refrech the session
func refrechCooki(w http.ResponseWriter) string {
	newToken := generateUID()
	setSessionCookie(w, newToken)
	return newToken
}
