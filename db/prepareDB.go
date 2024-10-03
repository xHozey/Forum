package forum

import (
	"database/sql"
	"fmt"
)

func PrepareDataBase(db *sql.DB) {
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
