package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/parser/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/router"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	contractRepository "github.com/Limechain/HCS-Integration-Node/app/domain/contract/repository"
	contractService "github.com/Limechain/HCS-Integration-Node/app/domain/contract/service"
	proposalRepository "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/repository"
	proposalService "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/service"
	poRepository "github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/repository"
	poService "github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/service"
	rfpRepository "github.com/Limechain/HCS-Integration-Node/app/domain/rfp/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	apiRouter "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/router"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common/queue"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/dlt/hcs"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/messaging/libp2p"
	contractMongo "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/contract"
	proposalMongo "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/proposal"
	poMongo "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/purchase-order"
	rfpMongo "github.com/Limechain/HCS-Integration-Node/app/persistance/mongodb/rfp"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func setupP2PClient(
	prvKey ed25519.PrivateKey,
	hcsClient common.Messenger,
	rfpRepo rfpRepository.RFPRepository,
	proposalRepo proposalRepository.ProposalRepository,
	contractRepo contractRepository.ContractsRepository,
	cs *contractService.ContractService,
	por poRepository.PurchaseOrdersRepository,
	pos *poService.PurchaseOrderService) common.Messenger {

	listenPort := os.Getenv("P2P_PORT")
	listenIp := os.Getenv("P2P_IP")
	peerMultiAddr := os.Getenv("PEER_ADDRESS")

	p2pClient := libp2p.NewLibP2PClient(prvKey, listenIp, listenPort, peerMultiAddr)

	// TODO get some env variables
	// TODO add more handlers
	rfpHandler := handler.NewRFPHandler(rfpRepo)
	proposalHandler := handler.NewProposalHandler(proposalRepo)
	contractRequestHandler := handler.NewContractRequestHandler(contractRepo, cs, p2pClient)
	contractAcceptedHandler := handler.NewContractAcceptedHandler(contractRepo, cs, hcsClient)
	poRequestHandler := handler.NewPORequestHandler(por, pos, p2pClient)
	poAcceptedHandler := handler.NewPOAcceptedHandler(por, pos, hcsClient)

	var parser json.JSONBusinessMesssageParser

	r := router.NewBusinessMessageRouter(&parser)

	r.AddHandler(messages.P2PMessageTypeRFP, rfpHandler)
	r.AddHandler(messages.P2PMessageTypeProposal, proposalHandler)
	r.AddHandler(messages.P2PMessageTypeContractRequest, contractRequestHandler)
	r.AddHandler(messages.P2PMessageTypeContractAccepted, contractAcceptedHandler)
	r.AddHandler(messages.P2PMessageTypePORequest, poRequestHandler)
	r.AddHandler(messages.P2PMessageTypePOAccepted, poAcceptedHandler)

	p2pChannel := make(chan *common.Message)

	p2pQueue := queue.New(p2pChannel, r)

	p2pClient.Listen(p2pQueue)

	return p2pClient
}

func setupDLTClient(
	prvKey ed25519.PrivateKey,
	contractRepo contractRepository.ContractsRepository,
	cs *contractService.ContractService,
	por poRepository.PurchaseOrdersRepository,
	pos *poService.PurchaseOrderService) common.Messenger {

	shouldConnectToMainnet := (os.Getenv("HCS_MAINNET") == "true")
	hcsClientID := os.Getenv("HCS_CLIENT_ID")
	hcsMirrorNodeID := os.Getenv("HCS_MIRROR_NODE_ADDRESS")
	topicID := os.Getenv("HCS_TOPIC_ID")

	var parser json.JSONBusinessMesssageParser

	r := router.NewBusinessMessageRouter(&parser)

	contractHandler := handler.NewDLTContractHandler(contractRepo, cs)
	poHandler := handler.NewDLTPOHandler(por, pos)
	// TODO add handlers

	r.AddHandler(messages.DLTMessageTypeContract, contractHandler)
	r.AddHandler(messages.DLTMessageTypePO, poHandler)

	ch := make(chan *common.Message)

	q := queue.New(ch, r)

	hcsClient := hcs.NewHCSClient(hcsClientID, prvKey, hcsMirrorNodeID, topicID, shouldConnectToMainnet)

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

	pubKey, ok := prvKey.Public().(ed25519.PublicKey)
	if !ok {
		log.Errorf("Could not cast the public key to ed25519 public key")
	}

	log.Infof("Started node with public key: %s\n", hex.EncodeToString(pubKey))

	peerPubKey := getPeerPublicKey()

	client, db := connectToDb()

	defer client.Disconnect(context.Background())

	rfpRepo := rfpMongo.NewRFPRepository(db)
	proposalRepo := proposalMongo.NewProposalRepository(db)
	contractRepo := contractMongo.NewContractRepository(db)
	por := poMongo.NewPurchaseOrderRepository(db)

	ps := proposalService.New()
	cs := contractService.New(prvKey, proposalRepo, ps, peerPubKey)
	pos := poService.New(prvKey, contractRepo, cs, peerPubKey)

	hcsClient := setupDLTClient(prvKey, contractRepo, cs, por, pos)

	defer hcsClient.Close()

	p2pClient := setupP2PClient(prvKey, hcsClient, rfpRepo, proposalRepo, contractRepo, cs, por, pos)

	defer p2pClient.Close()

	apiPort := os.Getenv("API_PORT")

	a := api.NewIntegrationNodeAPI()

	rfpApiService := apiservices.NewRFPService(rfpRepo, p2pClient)
	proposalApiService := apiservices.NewProposalService(proposalRepo, p2pClient)

	contractApiService := apiservices.NewContractService(contractRepo, cs, p2pClient)
	purchaseOrderApiService := apiservices.NewPurchaseOrderService(por, pos, p2pClient)

	a.AddRouter(fmt.Sprintf("/%s", apiRouter.RouteRFP), apiRouter.NewRFPRouter(rfpApiService))
	a.AddRouter(fmt.Sprintf("/%s", apiRouter.RouteProposal), apiRouter.NewProposalsRouter(proposalApiService))
	a.AddRouter(fmt.Sprintf("/%s", apiRouter.RouteContract), apiRouter.NewContractsRouter(contractApiService))
	a.AddRouter(fmt.Sprintf("/%s", apiRouter.RoutePO), apiRouter.NewPurchaseOrdersRouter(purchaseOrderApiService))

	if err := a.Start(apiPort); err != nil {
		panic(err)
	}

}
