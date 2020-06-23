package main

import (
	"context"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/parser/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/router"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/messaging/libp2p"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/queue"
	rfpPersistance "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/rfp"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"os/signal"
	"syscall"
)

const DefaultKeyPath = "./config/key.pem"

func main() {

	prvKey := getPrivateKey(DefaultKeyPath)

	messenger := libp2p.NewMessenger(prvKey)

	mongoConnString := os.Getenv("MONGODB_CONN_STR")
	mongoDatabaseName := os.Getenv("MONGODB_DBNAME")

	client, db := connectToDb(mongoConnString, mongoDatabaseName)

	defer client.Disconnect(context.Background())

	rfpRepo := rfpPersistance.NewRFPRepository(db)

	h := handler.NewRFPHandler(rfpRepo)

	var parser json.JSONBusinessMesssageParser

	router := router.NewBusinessMessageRouter(&parser)

	router.AddHandler("rfp", h)

	ch := make(chan *p2p.P2PMessage)

	q := queue.New(ch, router)

	messenger.Connect(q)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("[Ctrl + c] to shut down...")
	<-quit
	fmt.Println("Received exit signal, shutting down...")

}
