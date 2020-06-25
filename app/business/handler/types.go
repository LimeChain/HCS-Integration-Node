package handler

import (
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

const (
	P2PMessageTypeRFP              = "rfp"
	P2PMessageTypeProposal         = "proposal"
	P2PMessageTypeContractRequest  = "contractrequest"
	P2PMessageTypeContractAccepted = "contractaccepted"
	P2PMessageTypePO               = "po"
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
