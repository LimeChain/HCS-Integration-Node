package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func connectToDb() (*mongo.Client, *mongo.Database) {

	mongoConnString := os.Getenv("MONGODB_CONN_STR")
	mongoDatabaseName := os.Getenv("MONGODB_DBNAME")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConnString))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db := client.Database(mongoDatabaseName)

	return client, db
}
