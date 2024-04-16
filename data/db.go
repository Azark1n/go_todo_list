package data

import (
	"database/sql"
	"go_todo_list/config"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func OpenDbOrCreate() {
	if _, err := os.Stat(config.DbFile); os.IsNotExist(err) {
		log.Printf("The database file %s was not found, we are creating it...", config.DbFile)
		openDB()
		executeSQLFile(config.SchemaPath)
	} else {
		openDB()
	}
}

func executeSQLFile(filePath string) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("The file could not be read: %v", err)
	}

	_, err = Db.Exec(string(fileContent))
	if err != nil {
		log.Fatalf("Error when executing SQL from a file: %v", err)
	}
}

func openDB() {
	var err error
	Db, err = sql.Open("sqlite3", config.DbFile)
	if err != nil {
		log.Fatalf("Error opening the database: %v", err)
	}
}
