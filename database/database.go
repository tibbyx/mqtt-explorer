package database

import (
	"database/sql"
	"errors"
	//"github.com/mattn/go-sqlite3" // TODO: Uncomment it to use sqlite3.
)

func SetupDatabase() (*sql.DB, error) {
	err := errors.New("Not supported")
	return nil, err
}
