package handler

import (
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

const (
	BlockchainMessageTypeContract = "contract"
	BlockchainMessageTypePO       = "po"
)

type BusinessMessageParser interface {
	Parse(msg *common.Message) (*messages.BusinessMessage, error)
}

type BusinessLogicHandler interface {
	Handle(msg *common.Message) error
}
