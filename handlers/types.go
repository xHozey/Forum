package forum

import "database/sql"

// turn db to method to be used in handlers since we cant pass the db as arguments
type MyDB struct {
	MyData *sql.DB
}

// post informations
type Post struct {
	Id      int
	User    string
	Post    string
	Comment []string
	Like    int
	Deslike int
	Auth    bool
}

// struct with user name and if he's logged or not and All posts
type User struct {
	Auth     bool
	Username string
	Posts    []Post
}
