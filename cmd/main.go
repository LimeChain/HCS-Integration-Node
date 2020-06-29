package main

import (
	"context"
	"crypto/ed25519"
	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/parser/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/router"
	rfpRepository "github.com/Limechain/HCS-Integration-Node/app/domain/rfp/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	rfpRouter "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/router"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/blockchain/hcs"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common/queue"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/messaging/libp2p"
	rfpPersistance "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/rfp"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func setupP2PClient(prvKey ed25519.PrivateKey, rfpRepo rfpRepository.RFPRepository) common.Messenger {

	// TODO get some env variables
	// TODO add more handlers
	h := handler.NewRFPHandler(rfpRepo)

	var parser json.JSONBusinessMesssageParser

	r := router.NewBusinessMessageRouter(&parser)

	r.AddHandler(handler.P2PMessageTypeRFP, h)

	p2pChannel := make(chan *common.Message)

	p2pQueue := queue.New(p2pChannel, r)

	p2pClient := libp2p.NewLibP2PClient(prvKey)

	p2pClient.Listen(p2pQueue)

	return p2pClient
}

func setupBlockchainClient(prvKey ed25519.PrivateKey) common.Messenger {

	hcsClientID := os.Getenv("HCS_CLIENT_ID")
	hcsMirrorNodeID := os.Getenv("HCS_MIRROR_NODE_ADDRESS")
	topicID := os.Getenv("HCS_TOPIC_ID")

	var parser json.JSONBusinessMesssageParser

	r := router.NewBusinessMessageRouter(&parser)

	// TODO add handlers

	ch := make(chan *common.Message)

	q := queue.New(ch, r)

	hcsClient := hcs.NewHCSClient(hcsClientID, prvKey, hcsMirrorNodeID, topicID)

	err := hcsClient.Listen(q)
	if err != nil {
		panic(err)
	}

	return hcsClient

}

func main() {

	args := os.Args[1:]
	if len(args) > 0 {
		godotenv.Load(args[0])
	} else {
		godotenv.Load()
	}

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

	prvKey := getPrivateKey()

	client, db := connectToDb()

	defer client.Disconnect(context.Background())

	rfpRepo := rfpPersistance.NewRFPRepository(db)
	// TODO create more repos

	hcsClient := setupBlockchainClient(prvKey) // Pass it to the correct services instead of logging

	defer hcsClient.Close()

	p2pClient := setupP2PClient(prvKey, rfpRepo)

	defer p2pClient.Close()

	apiPort := os.Getenv("API_PORT")

	a := api.NewIntegrationNodeAPI()

	rfpService := apiservices.NewRFPService(rfpRepo)

	a.AddRouter("/rfp", rfpRouter.NewRFPRouter(rfpService))

	if err := a.Start(apiPort); err != nil {
		panic(err)
	}

}
