package main

import (
	"context"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/parser/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/router"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	rfpRouter "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/router"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/blockchain/hcs"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/messaging/libp2p"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/queue"
	rfpPersistance "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/rfp"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"os"
)

const DefaultKeyPath = "./config/key.pem"

func main() {

	logFilePath := os.Getenv("LOG_FILE")

	setupLogger()

	if len(logFilePath) > 0 {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Panic(err)
		}

		defer file.Close()

		setupFileLogger(file)
	}

	prvKey := getPrivateKey(DefaultKeyPath)

	hcsClientID := os.Getenv("HCS_CLIENT_ID")
	topicID := os.Getenv("HCS_TOPIC_ID")

	hcs.NewHCSClient(hcsClientID, prvKey, topicID) // TODO take the client and do something with it

	messenger := libp2p.NewMessenger(prvKey)

	mongoConnString := os.Getenv("MONGODB_CONN_STR")
	mongoDatabaseName := os.Getenv("MONGODB_DBNAME")

	client, db := connectToDb(mongoConnString, mongoDatabaseName)

	defer client.Disconnect(context.Background())

	rfpRepo := rfpPersistance.NewRFPRepository(db)

	h := handler.NewRFPHandler(rfpRepo)

	var parser json.JSONBusinessMesssageParser

	r := router.NewBusinessMessageRouter(&parser)

	r.AddHandler("rfp", h)

	ch := make(chan *p2p.P2PMessage)

	q := queue.New(ch, r)

	messenger.Connect(q)

	apiPort := os.Getenv("API_PORT")

	a := api.NewIntegrationNodeAPI()

	a.AddRouter("/rfp", rfpRouter.NewRFPRouter())

	if err := a.Start(apiPort); err != nil {
		panic(err)
	}

}
