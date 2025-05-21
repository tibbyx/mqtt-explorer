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
	queries := []string{
		`CREATE TABLE IF NOT EXISTS Broker (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Ip TEXT NOT NULL,
			Port INTEGER NOT NULL,
			CreationDate DATETIME NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS User (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			BrokerId INTEGER NOT NULL,
			ClientId INTEGER NOT NULL,
			Username TEXT NOT NULL,
			Password TEXT,
			Outsider BOOLEAN,
			CreationDate DATETIME NOT NULL,
			FOREIGN KEY(BrokerId) REFERENCES Broker(ID)
		);`,

		`CREATE TABLE IF NOT EXISTS Message (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserId INTEGER NOT NULL,
			TopicId INTEGER NOT NULL,
			QoS TINYINT,
			Date INTEGER,
			Message TEXT,
			FOREIGN KEY(UserId) REFERENCES User(ID),
			FOREIGN KEY(TopicId) REFERENCES Topic(ID)
		);`,
		`CREATE TABLE IF NOT EXISTS Topic (
             ID INTEGER PRIMARY KEY AUTOINCREMENT,
             UserId INTEGER NOT NULL,
             Subscribed BOOLEAN,
             Date INTEGER,
             Topic TEXT NOT NULL,
             FOREIGN KEY(UserId) REFERENCES User(ID)
         );`,

         `CREATE TABLE IF NOT EXISTS UserTopicFavourite (
             UserId INTEGER NOT NULL,
             TopicId INTEGER NOT NULL,
             Date INTEGER,
             PRIMARY KEY(UserId, TopicId),
             FOREIGN KEY(UserId) REFERENCES User(ID),
             FOREIGN KEY(TopicId) REFERENCES Topic(ID)
         );`,
	}

	for _, query := range queries {
		if _, err := con.Exec(query); err != nil {
			return fmt.Errorf("Skill issues\nErr: %s\n")
		}
    }


	return nil
}
