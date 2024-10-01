package handlers

import "database/sql"

type MyDB struct {
	MyData *sql.DB
}

type Post struct {
	User    string
	Post    string
	Comment string
	Like    int
	Deslike int
}
