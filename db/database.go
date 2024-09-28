package db

import (
	"database/sql"
	"fmt"
	"forum/types"
	"log"

	_ "github.com/mattn/go-sqlite3"
)



func (db types.DB) insertUser(mail, passwod string) error {
	query := `INSERT INTO login (mail, password) VALUE (?,?)`

	_, err := db.Data.Exec(query, mail, passwod )
	if err != nil {
		return err
	}
	fmt.Println("User add successfully")
	return nil
}

func EstableshDB()  {
	db , err := sql.Open("sqlite3", "./db/app.db")
	if err != nil {
		log.Fatal(err)
	}
}