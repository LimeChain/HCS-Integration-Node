package main

import (
	"database/sql"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler"
	"github.com/Limechain/HCS-Integration-Node/app/business/handler/parser/json"
	rfpHandler "github.com/Limechain/HCS-Integration-Node/app/business/handler/rfp"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p/queue"
	rfpPersistance "github.com/Limechain/HCS-Integration-Node/app/persistance/postgres/rfp"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var s *sql.DB

	repo := rfpPersistance.NewRFPRepository(s)

	h := rfpHandler.NewRFPHandler(repo)

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
