package database

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

const databaseName string = "mqtt-client-database.db"

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// |                | Polariusz | Created |
//
// # Description
// -  It opens a connection to a file `databaseName`. If it doesn't exist, it will be created.
// # Returns
// - DBcon to the database if everything went okay.
//
// # Author
// - Polariusz
func OpenDatabase() (*sql.DB, error) {
	con, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, fmt.Errorf("Error while opening database, perhaps there is a permission issue?\nErr: '%s'\n", err)
	}
	return con, nil
}

// | Date of change | By     | Comment |
// +----------------+--------+---------+
// | 2025-05-21     | Q-uock | Created |
//
// # Description
// - Creates tables in the connected to database connection.
//
// # Author
// - Q-uock
func SetupDatabase(con *sql.DB) error {
	tables := []string{
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

	for _, table := range tables {
		if _, err := con.Exec(table); err != nil {
			return fmt.Errorf("Skill issues\nErr: %s\n", err)
		}
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-21     | Polariusz | Created |
//
// # Struct to Table Mapping
//
// | Struct InsertBroker    | Table Broker          |
// +------------------------+-----------------------+
// |                        | ID Integer            |
// | Ip string              | Ip Text               |
// | Port int               | Port Integer          |
// |                        | CreationDate DateTime |
//
// # Used in
// - SelectBrokerList()
//
// # Author
// - Polariusz
type InsertBroker struct {
	Ip string
	Port int
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-21     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB        : It's a connection to the database that is used here to insert stuff in.
// - broker InsertBroker: The struct that will be written into table `Broker`, if it doesn't exists.
//
// # Description
// - The function shall insert the argument `broker` with the current date into table Broker from connected to database argument `con`.
// - The insertion shall only happen if the `broker` is not in the database.
//
// # Tables Affected
// - Broker
//   - INSERT
//   - SELECT (Subquery)
//
// # Returns
// - error if something bad happened. It can happen if the database that is connected does not have table Broker. In this case, please use functions `OpenDatabase()` and `SetupDatabase()` to set-up the database.
//
// # Author
// - Polariusz
func InsertNewBroker(con *sql.DB, broker InsertBroker) error {
	// I insert the arg broker while checking if it isn't in the database. If it is, the insertion will not happen.
	stmt, err := con.Prepare("INSERT INTO Broker(Ip, Port, CreationDate) SELECT ?, ?, ? WHERE NOT EXISTS(SELECT 1 FROM Broker WHERE Ip = ? AND Port = ?)");
	if err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(broker.Ip, broker.Port, time.Now(), broker.Ip, broker.Port); err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-21     | Polariusz | Created |
//
// # Struct to Table Mapping
//
// | Struct SelectBroker    | Table Broker          |
// +------------------------+-----------------------+
// | Id int                 | ID Integer            |
// | Ip string              | Ip Text               |
// | Port int               | Port Integer          |
// | CreationDate time.Time | CreationDate DateTime |
//
// # Used in
// - SelectBrokerList()
//
// # Author
// - Polariusz
type SelectBroker struct {
	Id int
	Ip string
	Port int
	CreationDate time.Time
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-21     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB: It's a connection to the database.
//
// # Description
// - Queries the table Broker from declared in function `SetupDatabase()` Tables.
// - It selects all rows.
//
// # Tables Affected
// - Broker
//   - SELECT
//
// # Returns
// - A list of struct `SelectBroker`
// - error can happen if the database that is connected does not have table Broker. In this case, please use functions `OpenDatabase()` and `SetupDatabase()` to set-up the database.
//
// # Author
// - Polariusz
func SelectBrokerList(con *sql.DB) ([]SelectBroker, error) {
	var brokerList []SelectBroker
	rows, err := con.Query("Select * from Broker")
	if err != nil {
		return nil, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	for rows.Next() {
		var Id int
		var Ip string
		var Port int
		var CreationDate time.Time

		rows.Scan(&Id, &Ip, &Port, &CreationDate)
		brokerList = append(brokerList, SelectBroker{Id, Ip, Port, CreationDate})
	}

	return brokerList, nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Struct to Table Mapping
//
// | Struct InsertTopic     | Table Topic           |
// +------------------------+-----------------------+
// |                        | ID INTEGER            |
// | UserId int             | UserId INTEGER        |
// | BrokerId int           | BrokerId INTEGER      |
// | Subscribed bool        | Subscribed BOOLEAN    |
// | Topic string           | Topic TEXT            |
// |                        | CreationDate DATETIME |
//
// # Used in
// - InsertNewTopic()
//
// # Author
// - Polariusz
type InsertTopic struct {
	UserId int
	BrokerId int
	Subscribed bool
	Topic string
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB       : It's a connection to the database that is used here to insert stuff in.
// - topic InsertTopic : Will be inserted into table `Topic`.
//
// # Description
// - The function shall insert the argument `topic` into table `Topic`.
// - The functino shall only insert unique argument `topic`.
//
// # Tables Affected
// - Topic
//   - INSERT
//     - SELECT (SUBQUERY)
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table Topic does not exists
//     - Use the `SetupDatabase()` function to set the database up before calling this function.
//   - Foreign Key issues
//
// # Author
// - Polariusz
func InsertNewTopic(con *sql.DB, topic InsertTopic) error {
	stmt, err := con.Prepare(`
		INSERT INTO Topic(UserId, BrokerId, Subscribed, Topic, CreationDate)
		Select ?, ?, ?, ?, ? WHERE NOT EXISTS(
			SELECT 1
			FROM Topic
			WHERE
				UserId = ?
			AND
				BrokerID = ?
			AND
				Topic = ?
		)
	`)
	if err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(topic.UserId, topic.BrokerId, topic.Subscribed, topic.Topic, time.Now()); err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database that is used here to insert stuff in.
// - topicId     : Unique Identifier of the Topic row
//
// # Description
// - The function shall update a row matched to argument `topicId` to mark the column `Subscribed` as true.
//
// # Tables Affected
// - Topic
//   - UPDATE
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table Topic does not exists
//     - Use the `SetupDatabase()` function to set the database up before calling this function.
//
// # Author
// - Polariusz
func SubscribeTopic(con *sql.DB, topicId int) error {
	stmt, err := con.Prepare(`
		UPDATE Topic
		SET Subscribed = true
		WHERE ID = ?
	`)
	if err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(topicId); err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database that is used here to insert stuff in.
// - topicId     : Unique Identifier of the Topic row
//
// # Description
// - The function shall update a row matched to argument `topicId` to mark the column `Subscribed` as false.
//
// # Tables Affected
// - Topic
//   - UPDATE
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table Topic does not exists
//     - Use the `SetupDatabase()` function to set the database up before calling this function.
//
// # Author
// - Polariusz
func UnsubscribeTopic(con *sql.DB, topicId int) error {
	stmt, err := con.Prepare(`
		UPDATE Topic
		SET Subscribed = false
		WHERE ID = ?
	`)
	if err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(topicId); err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB       : It's a connection to the database that is used here to insert stuff in.
// - topicId     : Unique Identifier of the Topic row
//
// # Description
// - The function shall remove a row from table `Topic` matched to argument `topicId`.
//
// # Tables Affected
// - Topic
//   - DELETE
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table Topic does not exists
//     - Use the `SetupDatabase()` function to set the database up before calling this function.
//
// # Author
// - Polariusz
func DeleteTopic(con *sql.DB, topicId int) error {
	stmt, err := con.Prepare(`
		DELETE FROM Topic
		WHERE ID = ?
	`)
	if err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(topicId); err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Struct to Table Mapping
//
// | Struct SelectTopic     | Table Topic           |
// +------------------------+-----------------------+
// | Id int                 | ID INTEGER            |
// | UserId int             | UserId INTEGER        |
// | BrokerId int           | BrokerId INTEGER      |
// | Subscribed bool        | Subscribed BOOLEAN    |
// | Topic string           | Topic TEXT            |
// | CreationDate time.Time | CreationDate DATETIME |
//
// # Used in
// - SelectSubscribedTopics()
// - SelectUnsubscribedTopics()
// - SelectTopicsByBrokerIdAndUserId()
//
// # Author
// - Polariusz
type SelectTopic struct {
	Id int
	UserId int
	BrokerId int
	Subscribed bool
	Topic string
	CreationDate time.Time
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database that is used here to insert stuff in.
// - brokerId    : Unique Identifier of table `Broker.ID`
// - userId      : Unique Identifier of table `User.ID`
//
// # Description
// - The function shall return an array of subscribed Topics matched with arguments `brokerId` and `userId`.
//
// # Tables Affected
// - Topic
//   - SELECT
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table Topic does not exists
//     - Use the `SetupDatabase()` function to set the database up before calling this function.
//
// # Author
// - Polariusz
func SelectSubscribedTopics(con *sql.DB, brokerId int, userId int) ([]SelectTopic, error) {
	var topicList []SelectTopic

	stmt, err := con.Prepare(`
		SELECT *
		FROM Topic
		WHERE
			BrokerId = ?
		AND
			UserID = ?
		AND
			Subscribed = true
	`)
	if err != nil {
		return nil, fmt.Errorf("Error while preparing the statement.\nErr: %s\n", err)
	}

	rows, err := stmt.Query(brokerId, userId)
	if err != nil {
		return nil, fmt.Errorf("Error while quering the statement.\nErr: %s\n", err)
	}

	for rows.Next() {
		var topic SelectTopic
		rows.Scan(&topic.Id, &topic.UserId, &topic.BrokerId, &topic.Subscribed, &topic.Topic, &topic.CreationDate)
		topicList = append(topicList, topic)
	}

	return topicList, nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database that is used here to insert stuff in.
// - brokerId    : Unique Identifier of table `Broker.ID`
// - userId      : Unique Identifier of table `User.ID`
//
// # Description
// - The function shall return an array of unsubscribed Topics matched with arguments `brokerId` and `userId`.
//
// # Tables Affected
// - Topic
//   - SELECT
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table Topic does not exists
//     - Use the `SetupDatabase()` function to set the database up before calling this function.
//
// # Author
// - Polariusz
func SelectUnsubscribedTopics(con *sql.DB, brokerId int, userId int) ([]SelectTopic, error) {
	var topicList []SelectTopic

	stmt, err := con.Prepare(`
		SELECT *
		FROM Topic
		WHERE
			BrokerId = ?
		AND
			UserID = ?
		AND
			Subscribed = false
	`)
	if err != nil {
		return nil, fmt.Errorf("Error while preparing the statement.\nErr: %s\n", err)
	}

	rows, err := stmt.Query(brokerId, userId)
	if err != nil {
		return nil, fmt.Errorf("Error while quering the statement.\nErr: %s\n", err)
	}

	for rows.Next() {
		var topic SelectTopic
		rows.Scan(&topic.Id, &topic.UserId, &topic.BrokerId, &topic.Subscribed, &topic.Topic, &topic.CreationDate)
		topicList = append(topicList, topic)
	}

	return topicList, nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database that is used here to insert stuff in.
// - brokerId    : Unique Identifier of table `Broker.ID`
// - userId      : Unique Identifier of table `User.ID`
//
// # Description
// - The function shall return an array of all known Topics matched with arguments `brokerId` and `userId`.
//
// # Tables Affected
// - Topic
//   - SELECT
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table Topic does not exists
//     - Use the `SetupDatabase()` function to set the database up before calling this function.
//
// # Author
// - Polariusz
func SelectTopicsByBrokerIdAndUserId(con *sql.DB, brokerId int, userId int) ([]SelectTopic, error) {
	var topicList []SelectTopic

	stmt, err := con.Prepare(`
		SELECT *
		FROM Topic
		WHERE
			BrokerId = ?
		AND
			UserID = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("Error while preparing the statement.\nErr: %s\n", err)
	}

	rows, err := stmt.Query(brokerId, userId)
	if err != nil {
		return nil, fmt.Errorf("Error while quering the statement.\nErr: %s\n", err)
	}

	for rows.Next() {
		var topic SelectTopic
		rows.Scan(&topic.Id, &topic.UserId, &topic.BrokerId, &topic.Subscribed, &topic.Topic, &topic.CreationDate)
		topicList = append(topicList, topic)
	}

	return topicList, nil
}
