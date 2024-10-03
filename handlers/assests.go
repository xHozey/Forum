package forum

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
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
