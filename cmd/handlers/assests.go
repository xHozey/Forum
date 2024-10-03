package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func (d *MyDB) insertLike(id int, user string) {
	var exists bool

	err := d.MyData.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user = ? AND post_id = ?)", user, id).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking like existence:", err)
		return
	}

	fmt.Println(exists)
	if !exists {
		_, err = d.MyData.Exec("INSERT INTO likes (user, post_id) VALUES (?, ?)", user, id)
		if err != nil {
			fmt.Println("Error inserting like:", err)
			return
		}

		_, err = d.MyData.Exec("UPDATE posts SET likes_count = likes_count + 1 WHERE id = ?", id)
		if err != nil {
			fmt.Println("Error updating likes count:", err)
		}
	} else {
		_, err = d.MyData.Exec("DELETE FROM likes WHERE user = ? AND post_id = ?", user, id)
		if err != nil {
			fmt.Println("Error removing like:", err)
		}

		_, err = d.MyData.Exec("UPDATE posts SET likes_count = likes_count - 1 WHERE id = ?", id)
		if err != nil {
			fmt.Println("Error updating likes count:", err)
		}
	}
}

func (d *MyDB) insertDeslike(id int, user string) {
	var exists bool
	err := d.MyData.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user = ? AND post_id = ?)", user, id).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !exists {
		_, err = d.MyData.Exec("UPDATE likes SET deslike = deslike + 1 WHERE id = ?", id)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err = d.MyData.Exec("UPDATE likes SET deslike = deslike - 1 WHERE id = ?", id)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (d *MyDB) getLikes(id int) (int, int) {
	var like, deslike int
	d.MyData.QueryRow("SELECT (like, deslike) FROM likes WHERE post_id = ?", id).Scan(&like, &deslike)

	return like, deslike
}

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
	stm, err := d.MyData.Query("SELECT comment FROM comments WHERE post_id = ? ORDER BY created_at DESC", postID)
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
	rows, err := d.MyData.Query("SELECT user, id, post FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.User, &post.Id, &post.Post)
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
