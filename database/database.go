package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Database represents the SQLite database connection and operations
type Database struct {
	db *sql.DB
}

// MQTT Message structure that maps to database table
type Message struct {
	ID        int64
	ClientID  string
	Topic     string
	Payload   string
	Timestamp string
}

// SetupDatabase initializes the SQLite database and creates tables if they don't exist
func SetupDatabase() (*Database, error) {
	// Find the project root (mqtt-explorer directory) regardless of where the application is run from
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}
	
	// Get the directory of the current file (database/database.go)
	currentDir := filepath.Dir(filename)
	
	// Go up one level to get the project root directory (mqtt-explorer)
	projectRoot := filepath.Dir(currentDir)
	
	// Ensure we have the correct directory
	if !strings.HasSuffix(projectRoot, "mqtt-explorer") {
		return nil, fmt.Errorf("could not find mqtt-explorer directory, current path: %s", projectRoot)
	}
	
	dataDir := filepath.Join(projectRoot, "data")
	
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Connect to SQLite database
	dbPath := filepath.Join(dataDir, "mqtt_explorer.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping the database to ensure connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create database instance
	database := &Database{db: db}

	// Initialize database schema
	if err := database.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	fmt.Println("Database setup completed successfully")
	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// initSchema creates necessary tables if they don't exist
func (d *Database) initSchema() error {
	// Create messages table
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			client_id TEXT NOT NULL,
			topic TEXT NOT NULL,
			payload TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create messages table: %w", err)
	}

	// Create topics table for tracking subscribed and known topics
	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS topics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			topic TEXT UNIQUE NOT NULL,
			is_subscribed BOOLEAN DEFAULT 0,
			first_seen DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create topics table: %w", err)
	}

	// Create connections table to track client connection history
	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS connections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			client_id TEXT NOT NULL,
			broker_host TEXT NOT NULL,
			broker_port TEXT NOT NULL,
			connected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			disconnected_at DATETIME
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create connections table: %w", err)
	}

	return nil
}

// SaveMessage saves an MQTT message to the database
func (d *Database) SaveMessage(clientID, topic, payload string) error {
	_, err := d.db.Exec(
		"INSERT INTO messages (client_id, topic, payload) VALUES (?, ?, ?)",
		clientID, topic, payload,
	)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}
	return nil
}

// GetMessagesByTopic retrieves all messages for a specific topic
func (d *Database) GetMessagesByTopic(topic string) ([]Message, error) {
	rows, err := d.db.Query(
		"SELECT id, client_id, topic, payload, timestamp FROM messages WHERE topic = ? ORDER BY timestamp DESC",
		topic,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ClientID, &msg.Topic, &msg.Payload, &msg.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return messages, nil
}

// AddOrUpdateTopic adds a new topic or updates an existing one
func (d *Database) AddOrUpdateTopic(topic string, isSubscribed bool) error {
	_, err := d.db.Exec(
		`INSERT INTO topics (topic, is_subscribed) 
		 VALUES (?, ?) 
		 ON CONFLICT(topic) DO UPDATE SET is_subscribed = ?`,
		topic, isSubscribed, isSubscribed,
	)
	if err != nil {
		return fmt.Errorf("failed to add/update topic: %w", err)
	}
	return nil
}

// GetAllTopics returns all known topics
func (d *Database) GetAllTopics() ([]string, error) {
	rows, err := d.db.Query("SELECT topic FROM topics ORDER BY first_seen DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query topics: %w", err)
	}
	defer rows.Close()

	var topics []string
	for rows.Next() {
		var topic string
		if err := rows.Scan(&topic); err != nil {
			return nil, fmt.Errorf("failed to scan topic row: %w", err)
		}
		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return topics, nil
}

// GetSubscribedTopics returns only subscribed topics
func (d *Database) GetSubscribedTopics() ([]string, error) {
	rows, err := d.db.Query("SELECT topic FROM topics WHERE is_subscribed = 1 ORDER BY first_seen DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query subscribed topics: %w", err)
	}
	defer rows.Close()

	var topics []string
	for rows.Next() {
		var topic string
		if err := rows.Scan(&topic); err != nil {
			return nil, fmt.Errorf("failed to scan topic row: %w", err)
		}
		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return topics, nil
}

// LogConnection records a new connection to the MQTT broker
func (d *Database) LogConnection(clientID, brokerHost, brokerPort string) (int64, error) {
	result, err := d.db.Exec(
		"INSERT INTO connections (client_id, broker_host, broker_port) VALUES (?, ?, ?)",
		clientID, brokerHost, brokerPort,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to log connection: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}
	
	return id, nil
}

// LogDisconnection records when a client disconnects from the MQTT broker
func (d *Database) LogDisconnection(connectionID int64) error {
	_, err := d.db.Exec(
		"UPDATE connections SET disconnected_at = CURRENT_TIMESTAMP WHERE id = ?",
		connectionID,
	)
	if err != nil {
		return fmt.Errorf("failed to log disconnection: %w", err)
	}
	return nil
}

// GetLastNMessages retrieves the last N messages from all topics
func (d *Database) GetLastNMessages(n int) ([]Message, error) {
	rows, err := d.db.Query(
		"SELECT id, client_id, topic, payload, timestamp FROM messages ORDER BY timestamp DESC LIMIT ?",
		n,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query last messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ClientID, &msg.Topic, &msg.Payload, &msg.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan message row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return messages, nil
}
