package database

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/go-kivik/couchdb/v3" // Import the CouchDB driver
	"github.com/go-kivik/kivik/v3"
)

// Client holds the CouchDB client connection (exported)
var Client *kivik.Client

// InitDB initializes the CouchDB client and creates the database if it doesn’t exist
func InitDB() error {
	// Use environment variables for configuration
	host := os.Getenv("COUCHDB_HOST")
	if host == "" {
		host = "202.178.125.77:5984" // Just the host and port, no protocol
	}
	user := os.Getenv("COUCHDB_USER")
	if user == "" {
		user = "admin"
	}
	pass := os.Getenv("COUCHDB_PASSWORD")
	if pass == "" {
		pass = "adminpw"
	}
	dbName := os.Getenv("COUCHDB_DATABASE")
	if dbName == "" {
		dbName = "ishopdb"
	}

	// Create connection string
	connString := fmt.Sprintf("http://%s:%s@%s", user, pass, host)

	// Initialize client
	var err error
	Client, err = kivik.New("couch", connString)
	if err != nil {
		log.Printf("Failed to connect to CouchDB: %v", err)
		return fmt.Errorf("failed to connect to CouchDB: %w", err)
	}

	fmt.Println("Database connected successfully")

	// Ensure the database exists
	ctx := context.Background()
	err = Client.CreateDB(ctx, dbName)
	if err != nil {
		// Check if the error is due to the database already existing (HTTP 412)
		if kivik.StatusCode(err) == 412 { // Precondition Failed (database exists)
			log.Printf("Database %s already exists, proceeding...", dbName)
		} else {
			log.Printf("Failed to create database %s: %v", dbName, err)
			return fmt.Errorf("failed to create database %s: %w", dbName, err)
		}
	}

	// Initialize views
	if err := initializeViews(Client.DB(ctx, dbName)); err != nil {
		log.Printf("Failed to initialize views: %v", err)
		return fmt.Errorf("failed to initialize views: %w", err)
	}

	return nil
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
	return nil
}
