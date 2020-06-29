package handler

import (
	"encoding/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type RFPHandler struct {
	rfpRepo repository.RFPRepository
}

func (h *RFPHandler) Handle(msg *common.Message) error {
	log.Infoln("Handling: ", string(msg.Msg))
	var rfpMsg messages.RFPMessage
	err := json.Unmarshal(msg.Msg, &rfpMsg)
	if err != nil {
		return err
	}
	rfpId, err := h.rfpRepo.Save(&rfpMsg.Data)
	if err != nil {
		return err
	}
	log.Infof("Saved rfp with id: %s\n", rfpId)
	return nil
}

func NewRFPHandler(repo repository.RFPRepository) *RFPHandler {
	return &RFPHandler{repo}
}
