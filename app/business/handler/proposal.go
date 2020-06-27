package handler

import (
	"encoding/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/proposal/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type ProposalHandler struct {
	proposalRepo repository.ProposalRepository
}

func (h *ProposalHandler) Handle(msg *common.Message) error {

	var proposalMsg messages.ProposalMessage
	err := json.Unmarshal(msg.Msg, &proposalMsg)
	if err != nil {
		return err
	}
	rfpId, err := h.proposalRepo.Save(&proposalMsg.Data)
	if err != nil {
		return err
	}
	log.Infof("Saved proposal with id: %s\n", rfpId)
	return nil
}

func NewProposalHandler(repo repository.ProposalRepository) *ProposalHandler {
	return &ProposalHandler{repo}
}
