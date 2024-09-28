package handlers

import "database/sql"

type MyDB struct {
	MyData *sql.DB
}
