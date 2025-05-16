package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}


func (d *Database) InitDatabase() error {
	var err error
	d.db, err = sql.Open("sqlite3", "./mqtt.db")
	if err != nil {
		return err
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS clients (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        adresse TEXT NOT NULL,
        port INTEGER NOT NULL,
        clientId TEXT NOT NULL
    );`


    _, err = d.db.Exec(createTable)
    return err
}

func (d *Database) SaveClientToDB(adresse string, port int, clientid string) error {
    _, err := d.db.Exec("INSERT INTO clients (adresse, port, clientId) VALUES (?, ?, ?)", adresse, port, clientId)
    return err
}

}
