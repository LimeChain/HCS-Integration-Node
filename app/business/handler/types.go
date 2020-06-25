package handler

import (
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

const P2PMessageTypeRFP = "rfp"
const P2PMessageTypeProposal = "proposal"
const P2PMessageTypeContractRequest = "contractrequest"
const P2PMessageTypeContractAccepted = "contractaccepted"
const P2PMessageTypePO = "po"

type BusinessMessageParser interface {
	Parse(msg *common.Message) (*messages.BusinessMessage, error)
}

type BusinessLogicHandler interface {
	Handle(msg *common.Message) error
}
