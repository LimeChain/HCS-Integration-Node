package main

import (
	"context"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/parser/json"
	rfpHandler "github.com/Limechain/HCS-Integration-Node/app/business/handler/rfp"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/queue"
	rfpPersistance "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/rfp"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func connectToDb(connString string) (*mongo.Client, *mongo.Database) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connString))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db := client.Database("hcs-integration-node")

	return client, db
}

func main() {

	mongoConnString := os.Getenv("MONGODB_CONN_STR")

	client, db := connectToDb(mongoConnString)

	defer client.Disconnect(context.Background())

	rfpRepo := rfpPersistance.NewRFPRepository(db)

	h := rfpHandler.NewRFPHandler(rfpRepo)

	var parser json.JSONBusinessMesssageParser

	router := handler.NewBusinessMessageRouter(&parser)

	router.AddHandler("rfp", h)

	ch := make(chan *p2p.P2PMessage)

	queue.New(ch, router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("[Ctrl + c] to shut down...")
	<-quit
	fmt.Println("Received exit signal, shutting down...")

}
