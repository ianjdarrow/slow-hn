package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func InitDB() *sqlx.DB {
	db = sqlx.MustConnect("sqlite3", "./db/slow-hn.db")
	err := db.Ping()
	if err != nil {
		fmt.Printf("Error establishing database connection: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Database connection established\n")
	CreateTables()
	return db
}

func CreateTables() {
	postSchema := `
    CREATE TABLE IF NOT EXISTS posts (
      by VARCHAR(255) ,
      id INTEGER PRIMARY KEY,
      score INTEGER,    
      time INTEGER,  
      title TEXT,  
      type VARCHAR(255),  
      url TEXT  
    );`
	scoreSchema := `
    CREATE TABLE IF NOT EXISTS scores (
      id INTEGER,
      score FLOAT,
      time INTEGER
    );`
	db.MustExec(postSchema)
	db.MustExec(scoreSchema)
}
