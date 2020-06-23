package handler

import (
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/p2p"
)

type BusinessMessageParser interface {
	Parse(msg *p2p.P2PMessage) (*messages.BusinessMessage, error)
}

type BusinessLogicHandler interface {
	Handle(msg *p2p.P2PMessage) error
}
