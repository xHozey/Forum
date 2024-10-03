package forum

import "database/sql"

type MyDB struct {
	MyData *sql.DB
}

type Post struct {
	Id      int
	User    string
	Post    string
	Comment []string
	Like    int
	Deslike int
	Auth    bool
}

type User struct {
	Auth     bool
	Username string
	Posts    []Post
}

