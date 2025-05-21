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
	clientConn *mongo.Client
	clientOnce sync.Once
)

func ConnectClient() *mongo.Client {
	clientOnce.Do(func() {
		if errENV := godotenv.Load(); errENV != nil {
			log.Fatal("Error in loading .env file")
		}

		db_username := os.Getenv("Mongo_username")
		db_password := os.Getenv("Mongo_password")
		db_template := os.Getenv("Mongo_URI")

		MongoDBURI := fmt.Sprintf(db_template, db_username, db_password)

		clientOpts := options.Client().
			ApplyURI(MongoDBURI).
			SetMinPoolSize(5).
			SetMaxPoolSize(50).
			SetMaxConnIdleTime(10 * time.Second).
			SetConnectTimeout(5 * time.Second)

		cli, err := mongo.Connect(clientOpts)
		if err != nil {
			log.Fatalf("Connection error: %v", err)
		}

		clientConn = cli

		fmt.Println("Successfully connected to MongoDB")
	})

	return clientConn
}
