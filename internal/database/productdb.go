package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-kivik/couchdb/v3" // Import the CouchDB driver
	"github.com/go-kivik/kivik/v3"
)

// Client holds the CouchDB client connection (exported)
var Client *kivik.Client

// Config holds the configuration for the CouchDB connection
type Config struct {
	Host     string
	Username string
	Password string
	Database string
	Port string
}

// InitDB initializes the CouchDB client and creates the database if it doesn’t exist
func InitDB() error {
	// Load configuration from environment variables
	cfg, err := loadConfig()
	if err != nil {
		log.Printf("Failed to load CouchDB configuration: %v", err)
		return fmt.Errorf("failed to load CouchDB configuration: %w", err)
	}

	// Create connection string
	connString := fmt.Sprintf("http://%s:%s@%s:%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port)

	// Initialize client with retry logic
	const maxRetries = 5
	const retryDelay = 2 * time.Second
	for retries := maxRetries; retries > 0; retries-- {
		Client, err = kivik.New("couch", connString)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to CouchDB (attempt %d/%d): %v", maxRetries-retries+1, maxRetries, err)
		if retries == 1 { // Last retry
			log.Printf("Exhausted retries connecting to CouchDB")
			return fmt.Errorf("failed to connect to CouchDB after %d attempts: %w", maxRetries, err)
		}
		time.Sleep(retryDelay)
	}

	log.Println("Database connected successfully")

	// Ensure the database exists
	ctx := context.Background()
	err = Client.CreateDB(ctx, cfg.Database)
	if err != nil {
		// Check if the error is due to the database already existing (HTTP 412)
		if kivik.StatusCode(err) == 412 { // Precondition Failed (database exists)
			log.Printf("Database %s already exists, proceeding...", cfg.Database)
		} else {
			log.Printf("Failed to create database %s: %v", cfg.Database, err)
			return fmt.Errorf("failed to create database %s: %w", cfg.Database, err)
		}
	}

	// Initialize views
	if err := initializeViews(Client.DB(ctx, cfg.Database)); err != nil {
		log.Printf("Failed to initialize views: %v", err)
		return fmt.Errorf("failed to initialize views: %w", err)
	}

	return nil
}

// loadConfig loads the CouchDB configuration from environment variables
func loadConfig() (Config, error) {
	cfg := Config{
		Host:     os.Getenv("COUCHDB_HOST"),
		Port:     os.Getenv("COUCHDB_PORT"), // Include port
		Username: os.Getenv("COUCHDB_USER"),
		Password: os.Getenv("COUCHDB_PASSWORD"),
		Database: os.Getenv("COUCHDB_DATABASE"),
	}

	// Validate required environment variables
	if cfg.Host == "" {
		return Config{}, fmt.Errorf("COUCHDB_HOST environment variable is required")
	}
	if cfg.Port == "" {
		return Config{}, fmt.Errorf("COUCHDB_PORT environment variable is required")
	}
	if cfg.Username == "" {
		return Config{}, fmt.Errorf("COUCHDB_USER environment variable is required")
	}
	if cfg.Password == "" {
		return Config{}, fmt.Errorf("COUCHDB_PASSWORD environment variable is required")
	}
	if cfg.Database == "" {
		return Config{}, fmt.Errorf("COUCHDB_DATABASE environment variable is required")
	}

	return cfg, nil
}


// GetDB returns a handle to the specified database
func GetDB(databaseName string) *kivik.DB {
	if Client == nil {
		log.Fatal("Database client not initialized. Call InitDB first.")
	}
	return Client.DB(context.Background(), databaseName)
}

// GetDBWithContext returns a handle to the specified database with a custom context
func GetDBWithContext(ctx context.Context, databaseName string) *kivik.DB {
	if Client == nil {
		log.Fatal("Database client not initialized. Call InitDB first.")
	}
	return Client.DB(ctx, databaseName)
}

// initializeViews sets up necessary CouchDB views
func initializeViews(db *kivik.DB) error {
	designDoc := map[string]interface{}{
		"_id": "_design/products",
		"views": map[string]interface{}{
			"by_name": map[string]string{
				"map": "function(doc) { if (doc.name) emit(doc.name, doc._id); }",
			},
		},
	}
	_, err := db.Put(context.Background(), "_design/products", designDoc)
	if err != nil {
		// Check if the error is due to a conflict (HTTP 409)
		if kivik.StatusCode(err) == 409 { // Conflict (design doc exists)
			log.Printf("View _design/products already exists, skipping...")
			return nil
		}
		return fmt.Errorf("failed to create view: %w", err)
	}
	log.Println("Successfully created view _design/products")
	return nil
}