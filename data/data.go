package data

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

var db *sql.DB

func init() {
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err = sql.Open("postgres", "dbname=myblog sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return
}
