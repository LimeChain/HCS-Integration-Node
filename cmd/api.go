package main

import (
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
)

func startAPI(port string) error {
	a := api.NewIntegrationNodeAPI()

	// TODO a.AddRouter()

	return a.Start(port)
}
