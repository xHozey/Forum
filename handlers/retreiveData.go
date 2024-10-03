package forum

import "fmt"

// get likes and deslikes count from posts table
func (d *MyDB) getLikes(id int) (int, int) {
	var like, deslike int
	d.MyData.QueryRow("SELECT likes_count, deslikes_count FROM posts WHERE id = ?", id).Scan(&like, &deslike)
	return like, deslike
}

// with post id we retreive all comments of this post since posts table is linked with comments table
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

// get all posts with its informations from the database
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
