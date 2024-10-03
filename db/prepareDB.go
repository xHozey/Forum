package forum

import (
	"database/sql"
	"fmt"
)

func PrepareDataBase(db *sql.DB) {
	// Creat login table that holds information about the user
	stm, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS login (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user TEXT NOT NULL,
	pass TEXT NOT NULL,
	uid TEXT
	)`)
	if err != nil {
		fmt.Println(err)
	}
	stm.Exec()
	// creat posts table that hold post with some informations
	stm, err = db.Prepare(`
	CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	likes_count INTEGER DEFAULT 0,
	deslikes_count INTEGER DEFAULT 0,
	user TEXT,
	post TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP)
	`)
	if err != nil {
		fmt.Println(err)
	}
	stm.Exec()
	// creat comments table that hold user comments linked with posts table id
	stm, err = db.Prepare(`
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
		likes_count INTEGER DEFAULT 0, 
		deslikes_count INTEGER DEFAULT 0,
        comment TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE)
`)
	if err != nil {
		fmt.Println(err)
	}
	stm.Exec()
	// creat like table to track who liked if he already liked and like again we must reset the like
	stm, err = db.Prepare(`CREATE TABLE IF NOT EXISTS likes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	user TEXT NOT NULL,
	UNIQUE (user, post_id),
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE)`)
	if err != nil {
		fmt.Println(err)
	}
	stm.Exec()
	// same as like table
	stm, err = db.Prepare(`CREATE TABLE IF NOT EXISTS deslikes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	user TEXT NOT NULL,
	UNIQUE (user, post_id),
	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE)`)
	if err != nil {
		fmt.Println(err)
	}
	stm.Exec()
}
