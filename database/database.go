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
// - InsertNewBroker()
// - SelectBrokerByIpAndPort()
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
// - SelectBrokerByIpAndPort()
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
// | 2025-05-22     | Polariusz | Created |
// 
// # Arguments
// - con *sql.DB        : It's a connection to the database that is used here to insert stuff in.
// - broker InsertBroker: Used to match a row from table Broker
// 
// # Description
// - The function shall return a matching to the argument `broker` full row from table Broker from connected to database argument `con`.
// - The function shall therefore allow for quering the Id of the table Broker if the Ip and Port are known.
// 
// # Tables Affected
// - Broker
//   - SELECT
// 
// # Returns
// - SelectBroker struct matched to the argument `broker`
// - error when a duplicate is present. This should never happen as long as the function `InsertNewBroker()` is used to insert the Brokers.
// 
// # Author
// - Polariusz
func SelectBrokerByIpAndPort(con *sql.DB, broker InsertBroker) (SelectBroker, error) {
	var fullBroker SelectBroker

	stmt, err := con.Prepare("SELECT * FROM BROKER WHERE Ip = ? AND Port = ?")
	if err != nil {
		return fullBroker, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows, err := stmt.Query(broker.Ip, broker.Port)
	if err != nil {
		return fullBroker, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows.Next()
	var Id int
	var Ip string
	var Port int
	var CreationDate time.Time
	rows.Scan(&Id, &Ip, &Port, &CreationDate)
	fullBroker = SelectBroker{Id, Ip, Port, CreationDate}

	if rows.Next() {
		// Duplicate detected!
		return fullBroker, fmt.Errorf("Error: Duplicate at table Broker! Args in: %s:%d", broker.Ip, broker.Port)
	}

	return fullBroker, nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-22     | Polariusz | Created |
//
// # Struct to Table Mapping
//
// | Struct InsertUser      | Table User            |
// +------------------------+-----------------------+
// |                        | ID INTEGER            |
// | BrokerId int           | BrokerId INTEGER      |
// | ClientId string        | ClientId TEXT         |
// | Username string        | Username TEXT         |
// | Password string        | Password TEXT         |
// | Outsider bool          | Outsider BOOLEAN      |
// |                        | CreationDate DATETIME |
//
// # Used in
// - InsertNewUser()
//
// # Author
// - Polariusz
type InsertUser struct {
	BrokerId int
	ClientId string
	Username string
	Password string
	Outsider bool
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-21     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB     : It's a connection to the database that is used here to insert stuff in.
// - user InsertUser : It's inserted to the table `User`
//
// # Description
// - The function shall insert the argument `user` with the current date into table User from connected to database argument `con`.
//
// # Tables Affected
// - User
//   - INSERT
//
// # Returns
// - Can return error if the con isn't connected or if it doesn't have table User. In this case, please use functions `OpenDatabase()` and `SetupDatabase()` to set-up the database.
//
// # Author
// - Polariusz
func InsertNewUser(con *sql.DB, user InsertUser) error {
	stmt, err := con.Prepare(`
		INSERT INTO User(BrokerId, ClientId, Username, Password, Outsider, CreationDate) VALUES (?, ?, ?, ?, ?, ?)
	`);
	if err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(user.BrokerId, user.ClientId, user.Username, user.Password, user.Outsider, time.Now()); err != nil {
		return fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	return nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-28     | Polariusz | Created |
//
// # Struct to Table Mapping
//
// | Struct SelectUser      | Table User            |
// +------------------------+-----------------------+
// | Id int                 | ID INTEGER            |
// | BrokerId int           | BrokerId INTEGER      |
// | ClientId string        | ClientId TEXT         |
// | Username string        | Username TEXT         |
// |                        | Password TEXT         |
// | Outsider bool          | Outsider BOOLEAN      |
// | CreationDate time.Time | CreationDate DATETIME |
//
// # Used in
// - SelectUserById()
// - SelectUsersByClientId()
//
// # Author
// - Polariusz
type SelectUser struct {
	Id int
	BrokerId int
	ClientId string
	Username string
	Outsider bool
	CreationDate time.Time
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-28     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB : It's a connection to the database that is used here to insert stuff in.
// - id int      : Unique Identifier of an User.
//
// # Description
// - The function shall query the database to return matched to argument `id` row from table User with a `SelectUser` struct.
//
// # Tables Affected
// - User
//   - SELECT
//
// # Returns
// - SelectUser struct matched to the argument `id`
// - error when:
//   - no match was found
//   - table User does not exist
//     - The Database was not prepared, run `SetupDatabase()` function before this.
//   - Skill issues
//
// # Author
// - Polariusz
func SelectUserById(con *sql.DB, id int) (SelectUser, error) {
	var user SelectUser

	stmt, err := con.Prepare(`
		SELECT ID, BrokerId, ClientId, UserrName, Outsider, CreationDate
		FROM User
		WHERE ID = ?
	`)
	if err != nil {
		return user, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return user, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if !rows.Next() {
		return user, fmt.Errorf("Table User matched to Id: %d yelded no results.\n", id)
	} else {
		rows.Scan(&user.Id, &user.BrokerId, &user.ClientId, &user.Username, &user.Outsider, &user.CreationDate)
	}

	return user, nil
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-28     | Polariusz | Created |
//
// # Arguments
// - con *sql.DB     : It's a connection to the database that is used here to insert stuff in.
// - ClientId string : Unique Identifier of an User.
//
// # Description
// - The function shall query the database to return matched to argument `clientId` row from table User with a `[]SelectUser` array of structs.
//
// # Tables Affected
// - User
//   - SELECT
//
// # Returns
// - []SelectUser array of structs matched to the argument `clientId`
//   - Why? Well, Because of the defined User and Broker Tables, it is possible to have same ClientIds to different Brokers.
//   - As the result of that it is possible for the clientId query to return multiple Users that were registered from different Broker IPs.
// - error when:
//   - no match was found
//   - table User does not exist
//     - The Database was not prepared, run `SetupDatabase()` function before this.
//   - Skill issues
//
// # Author
// - Polariusz
func SelectUsersByClientId(con *sql.DB, clientId string) ([]SelectUser, error) {
	var userList []SelectUser

	stmt, err := con.Prepare(`
		SELECT ID, BrokerId, ClientId, UserrName, Outsider, CreationDate
		FROM User
		WHERE ClientId = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows, err := stmt.Query(clientId)
	if err != nil {
		return nil, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	tableIsEmpty := true
	for rows.Next() {
		var user SelectUser
		tableIsEmpty = false

		rows.Scan(&user.Id, &user.BrokerId, &user.ClientId, &user.Username, &user.Outsider, &user.CreationDate)
		userList = append(userList, user)
	}
	if tableIsEmpty {
		return nil, fmt.Errorf("Table User matched to ClientId: %s yelded no results.\n", clientId)
	}

	return userList, nil
}
