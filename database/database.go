package database

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

const databaseName string = "mqtt-client-database.db"
const LIMIT_MESSAGES int = 500

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
// # Struct to Table Message
//
// | Struct InsertMessage   | Table Message         |
// +------------------------+-----------------------+
// |                        | ID INTEGER            |
// | UserId int             | UserId INTEGER        |
// | TopicId int            | TopicId INTEGER       |
// | BrokerId int           | BrokerId INTEGER      |
// | QoS int                | QoS TINYINT           |
// | Message string         | Message TEXT          |
// |                        | CreationDate DateTime |
//
// # Used in
// - InsertNewMessage()
//
// # Author
// - Polariusz
type InsertMessage struct {
	UserId int
	TopicId int
	BrokerId int
	QoS int
	Message string
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB           : It's a connection to the database that is used here to insert stuff in.
// - message InsertMessage : The struct that will be written into table `Message`.
//
// # Description
// - The function shall insert the argument `message` with the current date into table Message from connected to database argument `con`.
//
// # Tables Affected
// - Message
//   - INSERT
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table Message does not exist
//     - Run SetupDatabase() before this function.
//   - Foreign Key issues
//
// # Author
// - Polariusz
func InsertNewMessage(con *sql.DB, message InsertMessage) error {
	stmt, err := con.Prepare(`
		INSERT INTO Message(UserId, TopicId, BrokerId, QoS, Message, CreationDate)
		VALUES(?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(message.UserId, message.TopicId, message.BrokerId, message.QoS, message.Message, time.Now()); err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-29     | Polariusz | Created |
//
// # Struct to Table Message
//
// | Struct SelectMessage   | Table Message         |
// +------------------------+-----------------------+
// | Id int                 | ID INTEGER            |
// | UserId int             | UserId INTEGER        |
// | TopicId int            | TopicId INTEGER       |
// | BrokerId int           | BrokerId INTEGER      |
// | QoS int                | QoS TINYINT           |
// | Message string         | Message TEXT          |
// | CreationDate time.Time | CreationDate DateTime |
//
// # Used in
// - SelectMessagesByTopicIdAndBrokerId()
// - SelectMessagesByTopicIdBrokerIdAndIndex()
//
// # Author
// - Polariusz
type SelectMessage struct {
	Id int
	UserId int
	TopicId int
	BrokerId int
	QoS int
	Message string
	CreationDate time.Time
}

// | Date of change | By        | Comment                         |
// +----------------+-----------+---------------------------------+
// | 2025-05-29     | Polariusz | Created                         |
// | 2025-05-30     | Polariusz | Fixed references in rows.Scan() |
//
// # Arguments
// - con *sql.DB  : It's a connection to the database.
// - topicId int  : Unique Identifier of table Topic
// - brokerId int : Unique Identifier of table Broker
//
// # Description
// - Selects a list of messages from table Message matched to arguments `topicId` for messages in a Topic and `brokerId` for messages in a broker.
// - It selects all rows.
//
// # Tables Affected
// - Message
//   - SELECT
//
// # Returns
// - A list of struct `SelectMessage`
// - error when:
//   - Skill Issues
//   - Table Message does not exist
//     - Run SetupDatabase() before this function.
//
// # Author
// - Polariusz
func SelectMessagesByTopicIdAndBrokerId(con *sql.DB, topicId int, brokerId int) ([]SelectMessage, error) {
	var selectMessageList []SelectMessage

	stmt, err := con.Prepare(`
		SELECT *
		FROM Message
		WHERE
		  TopicId = ?
		AND
		  BrokerId = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows, err := stmt.Query(topicId, brokerId)
	if err != nil {
		return nil, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	for rows.Next() {
		var selectMessage SelectMessage
		rows.Scan(&selectMessage.Id, &selectMessage.UserId, &selectMessage.TopicId, &selectMessage.BrokerId, &selectMessage.QoS, &selectMessage.Message, &selectMessage.CreationDate)
		selectMessageList = append(selectMessageList, selectMessage)
	}

	return selectMessageList, nil
}

// | Date of change | By        | Comment                                                                                    |
// +----------------+-----------+--------------------------------------------------------------------------------------------+
// | 2025-05-29     | Polariusz | Created                                                                                    |
// | 2025-05-30     | Polariusz | Fixed references in rows.Scan() and changed the statement to use the ROW_NUMBER() function |
//
// # Arguments
// - con *sql.DB  : It's a connection to the database.
// - topicId int  : Unique Identifier of table Topic
// - brokerId int : Unique Identifier of table Broker
// - index int    : Select from `LIMIT_MESSAGES*index` to `LIMIT_MESSAGES*(1+index)` messages.
//
// # Description
// - The function shall select matched to arguments `` for matching to Topic, `` for matching to Broker and `` for limiting messages Messages from table `Message` by a connected to `con` Database.
// - It selects up to `LIMIT_MESSAGES` Messages
//
// # Tables Affected
// - Message
//   - SELECT
//
// # Returns
// - A list of struct `SelectMessage`
// - error when:
//   - Skill Issues
//   - Table Message does not exist
//     - Run SetupDatabase() before this function.
//
// # Author
// - Polariusz
func SelectMessagesByTopicIdBrokerIdAndIndex(con *sql.DB, topicId int, brokerId int, index int) ([]SelectMessage, error) {
	var selectMessageList []SelectMessage
	stmtStr := `
		SELECT ID, UserId, TopicId, BrokerId, QoS, Message, CreationDate
		FROM (
			ROW_NUMBER() OVER(ORDER BY ID) as RowCnt, ID, UserId, TopicId, BrokerId, QoS, Message, CreationDate
			FROM MESSAGE
			WHERE
				TopicId = ?
			AND
				BrokerId = ?
		) MsgWithCnt
		WHERE
			RowCnt > ? * ?
		AND
			RowCnt <= (1+?) * ?
	`

	stmt, err := con.Prepare(stmtStr)
	if err != nil {
		return nil, fmt.Errorf("Error while preparing the statement!\nStatement:\n%s\nErr: %s\n", stmtStr, err)
	}

	rows, err := stmt.Query(topicId, brokerId, index, LIMIT_MESSAGES, index, LIMIT_MESSAGES)
	if err != nil {
		return nil, fmt.Errorf("Error while querying the statement!\nStatement:\n%s\nErr: %s\n", stmtStr, err)
	}

	for rows.Next() {
		var selectMessage SelectMessage
		rows.Scan(&selectMessage.Id, &selectMessage.UserId, &selectMessage.TopicId, &selectMessage.BrokerId, &selectMessage.QoS, &selectMessage.Message, &selectMessage.CreationDate)
		selectMessageList = append(selectMessageList, selectMessage)
	}

	return selectMessageList, nil
}
