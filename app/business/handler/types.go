package handler

import (
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
)

type BusinessMessage struct {
	Type string `json:"type"`
}

type BusinessMessageParser interface {
	Parse(msg *p2p.P2PMessage) (*BusinessMessage, error)
}

type BusinessLogicHandler interface {
	Handle(msg *p2p.P2PMessage) error
}
