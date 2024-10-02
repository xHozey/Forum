package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func (d *MyDB) insertComment(cmn string, postID int) error {
	var exists bool
	err := d.MyData.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = d.MyData.Exec("INSERT INTO comments (post_id, comment) VALUES (?, ?)", postID, cmn)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("post with id %d does not exist", postID)
	}

	return nil
}

func (d *MyDB) getComment(postID int) []string {
	stm, err := d.MyData.Query("SELECT comment FROM comments WHERE post_id = ?", postID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer stm.Close()

	var comments []string
	for stm.Next() {
		var comment string
		err := stm.Scan(&comment)
		if err != nil {
			fmt.Println(err)
			continue
		}
		comments = append(comments, comment)
	}

	if err = stm.Err(); err != nil {
		fmt.Println(err)
	}

	return comments
}

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

func (d *MyDB) getPosts() ([]Post, error) {
	rows, err := d.MyData.Query("SELECT user, id, likes, deslikes, post FROM posts")
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
