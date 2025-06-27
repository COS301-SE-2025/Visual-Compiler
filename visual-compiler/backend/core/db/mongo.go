// Package path implements routines for the MongoDB Connection Singleton
package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/joho/godotenv"
)

var (
	client_conn *mongo.Client
	client_once sync.Once
)

// Name: ConnectClient
//
// Parameters: None
//
// Return: Client instance
//
// This function creates the MongoDB connection as a Singleton. Decouples the connection logic from the CRUD operations
func ConnectClient() *mongo.Client {
	client_once.Do(func() {
		if _, err := os.Stat(".env"); err == nil {
			if err_env := godotenv.Load(); err_env != nil {
				log.Println("Warning. Failed to load .env file")
			} else {
				fmt.Println(".env file was successfully loaded")
			}
		} else {
			fmt.Println(".env file not found. Assuming environmental variables")
		}

		db_username := os.Getenv("Mongo_username")
		db_password := os.Getenv("Mongo_password")
		db_template := os.Getenv("Mongo_URI")

		if db_username == "" || db_password == "" || db_template == "" {
			log.Fatal("Missing variables")
		}

		mongo_db_uri := fmt.Sprintf(db_template, db_username, db_password)

		client_opts := options.Client().
			ApplyURI(mongo_db_uri).
			SetMinPoolSize(5).
			SetMaxPoolSize(50).
			SetMaxConnIdleTime(10 * time.Second).
			SetConnectTimeout(5 * time.Second)

		cli, err := mongo.Connect(client_opts)
		if err != nil {
			log.Fatalf("Connection error: %v", err)
		}

		client_conn = cli

		fmt.Println("Successfully connected to MongoDB")
	})

	return client_conn
}
