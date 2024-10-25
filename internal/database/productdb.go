package database

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
)

var Client *kivik.Client

func InitDB() {

	// creates a new client object to connect to CouchDB
	client, err := kivik.New("couch", "http://admin:password@localhost:5984/")

	if err != nil {
		log.Fatal("Filed to connnect to CouchDB: ", err)
	}

	fmt.Println("Database connected succesfully")

	// Assign the client to the global variable
	Client = client

}

func GetDB(databaseName string) *kivik.DB {
	return Client.DB(context.TODO(), databaseName)
}
