package db

func Migrations() []string {
	Users := ` CREATE TABLE IF NOT EXISTS  users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);`

	Posts := `CREATE TABLE IF NOT EXISTS  posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    images TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);`

	Comments := `CREATE TABLE IF NOT EXISTS  comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE, -- Cascades delete on post removal
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Cascades delete on user removal
);`

	Categories := `CREATE TABLE  IF NOT EXISTS  categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);`

	PostCategories := `CREATE TABLE IF NOT EXISTS  post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE, -- Cascades delete on post removal
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE, -- Cascades delete on category removal
    PRIMARY KEY (post_id, category_id)
);`

	LikesDislikes := `CREATE TABLE IF NOT EXISTS  comment_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    is_like BOOLEAN NOT NULL,  -- TRUE for like, FALSE for dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE, -- Cascades delete on comment removal
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Cascades delete on user removal
);`

	Sessions := `CREATE TABLE IF NOT EXISTS  sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_token TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);`
	var slice []string
	slice = append(slice, Users)
	slice = append(slice, Posts)
	slice = append(slice, Comments)
	slice = append(slice, Categories)
	slice = append(slice, PostCategories)
	slice = append(slice, LikesDislikes)
	slice = append(slice, Sessions)


	return slice
}