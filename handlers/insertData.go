package forum

import "fmt"

func (d *MyDB) insertLike(id int, user string) {
	var exists bool

	err := d.MyData.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE user = ? AND post_id = ?)", user, id).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking like existence:", err)
		return
	}
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

	err := d.MyData.QueryRow("SELECT EXISTS(SELECT 1 FROM deslikes WHERE user = ? AND post_id = ?)", user, id).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking like existence:", err)
		return
	}
	if !exists {
		_, err = d.MyData.Exec("INSERT INTO deslikes (user, post_id) VALUES (?, ?)", user, id)
		if err != nil {
			fmt.Println("Error inserting like:", err)
			return
		}

		_, err = d.MyData.Exec("UPDATE posts SET deslikes_count = deslikes_count + 1 WHERE id = ?", id)
		if err != nil {
			fmt.Println("Error updating likes count:", err)
		}
	} else {
		_, err = d.MyData.Exec("DELETE FROM deslikes WHERE user = ? AND post_id = ?", user, id)
		if err != nil {
			fmt.Println("Error removing like:", err)
		}

		_, err = d.MyData.Exec("UPDATE posts SET deslikes_count = deslikes_count - 1 WHERE id = ?", id)
		if err != nil {
			fmt.Println("Error updating likes count:", err)
		}
	}
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
