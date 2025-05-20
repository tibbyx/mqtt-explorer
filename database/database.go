package database

import (
	"database/sql"
	"fmt"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

const databaseName string = "mqtt-client-database.db"

func OpenDatabase() (*sql.DB, error) {
	con, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while opening database, perhaps there is a permission issue?\nErr: '%s'\n", err))
	}
	return con, nil
}

func SetupDatabase(con *sql.DB) error {
	_, err := con.Exec("CREATE TABLE IF NOT EXISTS Broker(ID INTEGER PRIMARY KEY AUTOINCREMENT,Ip string NOT NULL,Port INTEGER NOT NULL, CreationDate DATETIME NOT NULL)")
	if err != nil {
		return errors.New(fmt.Sprintf("Skill issues\nErr: %s\n", err))
	}
	// continue...
	return nil
}
