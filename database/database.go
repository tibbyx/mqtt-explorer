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
// | 2025-05-30     | Polariusz | Created |
//
// # Struct to Table Mapping
//
// | Struct SelectFavTopic  | Table UserTopicFavourite | Table User | Table Topic |
// +------------------------+--------------------------+------------+-------------+
// | Id int                 | ID INTEGER               |            |             |
// | UserId int             | UserId INTEGER           | ID INTEGER |             |
// | TopicId int            | TopicId INTEGER          |            | ID INTEGER  |
// | Topic string           |                          |            | Topic TEXT  |
// | CreationDate time.Time | CreationDate DATETIME    |            |             |
//
// # Used in
// - SelectFavouriteTopicsByUserId()
//
// # Author
// - Polariusz
type SelectFavTopic struct {
	Id int
	UserId int
	TopicId int
	Topic string
	CreationDate time.Time
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-30     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database.
// - userId int  : [User].[ID]
//
// # Description
// - The function shall return a list of favourite topics matched with argument `userId` with a `SelectFavTopic` struct array.
//
// # Tables Affected
// - UserTopicFavourite
//   - SELECT
//
// # Returns
// - An array of SelectFavTopic structs that match the argument `userId`
// - error when:
//   - Skill Issues
//   - Table UserTopicFavourite or Table Topic does not exist
//     - Please run the function `SetupDatabase()` before using the database.
//
// # Author
// - Polariusz
func SelectFavouriteTopicsByUserId(con *sql.DB, userId int) ([]SelectFavTopic, error) {
	var favTopicList []SelectFavTopic

	stmtStr := `
		SELECT utf.Id, utf.UserId, utf.TopicId, t.Topic, utf.CreationDate
		FROM UserTopicFavourite utf
		INNER JOIN Topic t
		ON t.ID = utf.TopicId
		WHERE utf.UserId = ?
	`
	stmt, err := con.Prepare(stmtStr)
	if err != nil {
		return nil, fmt.Errorf("Error while preparing the statement!\nStatement:\n%s\nErr: %s\n", stmtStr, err)
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying the database!\nStatement:\n%s\nErr: %s\n", stmtStr, err)
	}

	for rows.Next() {
		var favTopic SelectFavTopic
		rows.Scan(&favTopic.Id, &favTopic.UserId, &favTopic.TopicId, &favTopic.Topic, &favTopic.CreationDate)
		favTopicList = append(favTopicList, favTopic)
	}

	return favTopicList, nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-30     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database.
// - userId int  : [User].[ID]
// - topicId int : [Topic].[ID]
//
// # Description
// - The function shall insert the arguments `userId` and `topicId` into the table `UserTopicFavourite`.
//
// # Tables Affected
// - UserTopicFavourite
//   - INSERT
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table UserTopicFavourite does not exist
//     - Please run the function `SetupDatabase()` before using the database.
//
// # Author
// - Polariusz
func InsertFavouriteTopic(con *sql.DB, userId int, topicId int) error {
	stmtStr := `
		INSERT INTO UserTopicFavourite(UserId, TopicId, CreationDate)
		VALUES(?, ?, ?)
	`

	stmt, err := con.Prepare(stmtStr)
	if err != nil {
		return fmt.Errorf("Error while preparing the statement!\nStatement:\n%s\nErr: %s\n", stmtStr, err)
	}

	if _, err := stmt.Exec(userId); err != nil {
		return fmt.Errorf("Error while executing the statement!\nStatement:\n%s\nErr: %s\n", stmtStr, err)
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-30     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database.
// - userId int  : [User].[ID]
// - topicId int : [Topic].[ID]
//
// # Description
// - The function shall delete rows matched with arguments `userId` and `topicId` from the table `UserTopicFavourite`.
//
// # Tables Affected
// - UserTopicFavourite
//   - DELETE
//
// # Returns
// - error when:
//   - Skill Issues
//   - Table UserTopicFavourite does not exist
//     - Please run the function `SetupDatabase()` before using the database.
//
// # Author
// - Polariusz
func DeleteFavouriteTopic(con *sql.DB, userId int, topicId int) error {
	stmtStr := `
		DELETE FROM UserTopicFavourite
		WHERE userId = ?
		AND topicId = ?
	`

	stmt, err := con.Prepare(stmtStr)
	if err != nil {
		return fmt.Errorf("Error while preparing the statement!\nStatement:\n%s\nErr: %s\n", stmtStr, err)
	}

	if _, err := stmt.Exec(userId); err != nil {
		return fmt.Errorf("Error while executing the statement!\nStatement:\n%s\nErr: %s\n", stmtStr, err)
	}

	return nil
}
