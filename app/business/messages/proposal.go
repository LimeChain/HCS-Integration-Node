package messages

import "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/model"

type ProposalMessage struct {
	BusinessMessage
	Data model.Proposal `json:"data"`
}

func CreateProposalMessage(proposal *model.Proposal) *ProposalMessage {
	return &ProposalMessage{BusinessMessage: BusinessMessage{Type: "proposal"}, Data: *proposal}
}
