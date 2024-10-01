package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)
const file string = "Forum.db"

func GetDB() (*sql.DB, error) {
	DB, err := sql.Open("sqlite3", file)
	if err != nil {
		fmt.Println(err)
	}
	// defer DB.Close()
	return DB, nil
}

func BaseDonne() {
	// err := DB.Ping()
	// if err != nil {
	// 	log.Fatal("Cannot connect to database:", err)
	// }
	DB, err := GetDB()
	defer DB.Close()
	query := Migrations()
	for _, v := range query {
		_, err = DB.Exec(v)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database initialized.")
} // Create the items table if it doesn't exist
