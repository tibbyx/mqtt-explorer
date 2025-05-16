package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./mqtt.db")
	if err != nil {
		return err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		topic TEXT NOT NULL,
		payload TEXT NOT NULL,
		timestamp INTEGER NOT NULL
	);`

	_, err = db.Exec(createTable)
	return err
}

func SaveMessageToDB(topic string, payload string, timestamp int64) error {
	_, err := db.Exec("INSERT INTO messages (topic, payload, timestamp) VALUES (?, ?, ?)", topic, payload, timestamp)
	return err
}
