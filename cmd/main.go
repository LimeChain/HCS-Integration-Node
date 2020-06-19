package main

import (
	"context"
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

	var parser json.JSONBusinessMesssageParser

	router := handler.NewBusinessMessageRouter(&parser)

	var s *sql.DB

	repo := rfpPersistance.NewRFPRepository(s)

	fmt.Println(repo)

	h := rfpHandler.NewRFPHandler(repo)

	fmt.Println(h)

	router.AddHandler("rfp", h)

	ch := make(chan *p2p.P2PMessage)

	q := queue.New(ch, router)

	fmt.Println(q)
	q.Receive(&p2p.P2PMessage{[]byte("rfp"), context.Background()})

	exch := make(chan os.Signal, 1)
	signal.Notify(exch, syscall.SIGINT, syscall.SIGTERM)
	<-exch
	fmt.Println("Received signal, shutting down...")

}
