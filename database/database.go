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
			ClientId TEXT NOT NULL,
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
			BrokerId INTEGER NOT NULL,
			QoS TINYINT,
			Message TEXT,
			CreationDate DATETIME,
			FOREIGN KEY(UserId) REFERENCES User(ID),
			FOREIGN KEY(TopicId) REFERENCES Topic(ID),
			FOREIGN KEY(BrokerId) REFERENCES Broker(ID)
		);`,

		`CREATE TABLE IF NOT EXISTS Topic (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserId INTEGER NOT NULL,
			BrokerId INTEGER NOT NULL,
			Subscribed BOOLEAN,
			Topic TEXT NOT NULL,
			CreationDate DATETIME,
			FOREIGN KEY(UserId) REFERENCES User(ID),
			FOREIGN KEY(BrokerId) REFERENCES Broker(ID)
		);`,

		`CREATE TABLE IF NOT EXISTS UserTopicFavourite (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserId INTEGER NOT NULL,
			TopicId INTEGER NOT NULL,
			CreationDate DATETIME,
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

/*                                       +--------+                                       */
/* --------------------------------------| BROKER |-------------------------------------- */
/*                                       +--------+                                       */

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
//
// # Author
// - Polariusz
type InsertBroker struct {
	Ip string
	Port int
}

// | Date of change | By        | Comment             |
// +----------------+-----------+---------------------+
// | 2025-05-21     | Polariusz | Created             |
// | 2025-06-04     | Polariusz | Added the ID return |
//
// # Arguments
// - con *sql.DB        : It's a connection to the database that is used here to insert stuff in.
// - broker InsertBroker: The struct that will be written into table `Broker`, if it doesn't exists.
//
// # Description
// - The function shall insert the argument `broker` with the current date into table Broker from connected to database argument `con`.
// - The insertion shall only happen if the `broker` is not in the database.
// - The function shall return the ID of the inserted broker.
//
// # Tables Affected
// - Broker
//   - INSERT
//   - SELECT
//
// # Returns
// - int: it's the ID from table Broker that match the arguments `broker`. It will be -1 if an error has accured.
// - error if something bad happened. It can happen if the database that is connected does not have table Broker. In this case, please use functions `OpenDatabase()` and `SetupDatabase()` to set-up the database.
//
// # Author
// - Polariusz
func InsertNewBroker(con *sql.DB, broker InsertBroker) (int, error) {
	// I insert the arg broker while checking if it isn't in the database. If it is, the insertion will not happen.
	stmt, err := con.Prepare("INSERT INTO Broker(Ip, Port, CreationDate) SELECT ?, ?, ? WHERE NOT EXISTS(SELECT 1 FROM Broker WHERE Ip = ? AND Port = ?)");
	if err != nil {
		return -1, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(broker.Ip, broker.Port, time.Now(), broker.Ip, broker.Port); err != nil {
		return -1, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	// I want to get the ID of it.
	stmt, err = con.Prepare("SELECT ID from Broker WHERE Ip = ? AND Port = ?");
	if err != nil {
		return -1, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows, err := stmt.Query(broker.Ip, broker.Port)
	if err != nil {
		return -1, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows.Next()
	var ID int
	rows.Scan(&ID);

	return ID, nil
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

// | Date of change | By        | Comment                                |
// +----------------+-----------+----------------------------------------+
// | 2025-05-22     | Polariusz | Created                                |
// | 2025-06-02     | Polariusz | Changed the arguments for the function |
//
// # Arguments
// - con *sql.DB : It's a connection to the database that is used here to insert stuff in.
// - ip string   : Used to match to [Broker].[Ip]
// - port int    : Used to match to [Broker].[Port]
//
// # Description
// - The function shall return a matching to the arguments `ip` and `port` full row from table Broker from connected to database argument `con`.
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
func SelectBrokerByIpAndPort(con *sql.DB, ip string, port int) (SelectBroker, error) {
	var fullBroker SelectBroker

	stmt, err := con.Prepare("SELECT * FROM BROKER WHERE Ip = ? AND Port = ?")
	if err != nil {
		return fullBroker, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows, err := stmt.Query(ip, port)
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
		return fullBroker, fmt.Errorf("Error: Duplicate at table Broker! Args in: %s:%d", ip, port)
	}

	return fullBroker, nil
}

/*                                       +------+                                       */
/* --------------------------------------| USER |-------------------------------------- */
/*                                       +------+                                       */

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

// | Date of change | By        | Comment             |
// +----------------+-----------+---------------------+
// | 2025-05-22     | Polariusz | Created             |
// | 2025-06-04     | Polariusz | Added the ID return |
//
// # Arguments
// - con *sql.DB     : It's a connection to the database that is used here to insert stuff in.
// - user InsertUser : It's inserted to the table `User`
//
// # Description
// - The function shall insert the argument `user` with the current date into table User from connected to database argument `con`.
// - The function shall return the ID of the table User that match the inserted argument `user`.
//
// # Tables Affected
// - User
//   - INSERT
//   - SELECT
//
// # Returns
// - int: it's the ID from table User that match the arguments `user`. It will be -1 if an error has accured.
// - Can return error if the con isn't connected or if it doesn't have table User. In this case, please use functions `OpenDatabase()` and `SetupDatabase()` to set-up the database.
//
// # Author
// - Polariusz
func InsertNewUser(con *sql.DB, user InsertUser) (int, error) {
	stmt, err := con.Prepare(`
		INSERT INTO User(BrokerId, ClientId, Username, Password, Outsider, CreationDate) VALUES (?, ?, ?, ?, ?, ?)
	`);

	if err != nil {
		return -1, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	if _, err := stmt.Exec(user.BrokerId, user.ClientId, user.Username, user.Password, user.Outsider, time.Now()); err != nil {
		return -1, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	stmt, err := con.Prepare("SELECT ID FROM User WHERE BrokerId = ? AND ClientID = ? AND Username = ? AND Password = ? AND Outsider = ?")
	if err != nil {
		return -1, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows, err := con.Query(user.BrokerId, user.ClientId, user.Username, user.Password, user.Outsider)
	if err != nil {
		return -1, fmt.Errorf("Skill issues\nErr: %s\n", err)
	}

	rows.Next()
	var userId int
	rows.Scan(&userId)

	return userId, nil
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

/*                                       +-------+                                       */
/* --------------------------------------| TOPIC |-------------------------------------- */
/*                                       +-------+                                       */

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
// - con *sql.DB : It's a connection to the database that is used here to insert stuff in.
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

/*                                       +---------+                                       */
/* --------------------------------------| MESSAGE |-------------------------------------- */
/*                                       +---------+                                       */

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
// | 2025-06-02     | Polariusz | added missing arguments under the description documentation of the function                |
//
// # Arguments
// - con *sql.DB  : It's a connection to the database.
// - topicId int  : Unique Identifier of table Topic
// - brokerId int : Unique Identifier of table Broker
// - index int    : Select from `LIMIT_MESSAGES*index` to `LIMIT_MESSAGES*(1+index)` messages.
//
// # Description
// - The function shall select matched to arguments `topicId` for matching to Topic, `brokerId` for matching to Broker and `index` for limiting messages Messages from table `Message` by a connected to `con` Database.
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

/*                                       +----------+                                       */
/* --------------------------------------| FAVTOPIC |-------------------------------------- */
/*                                       +----------+                                       */

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
