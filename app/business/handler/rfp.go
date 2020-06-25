package handler

import (
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type RFPHandler struct {
	rfpRepo repository.RFPRepository
}

func (h *RFPHandler) Handle(msg *common.Message) error {
	fmt.Println("Handling: ", string(msg.Msg))
	return nil
}

func NewRFPHandler(repo repository.RFPRepository) *RFPHandler {
	return &RFPHandler{repo}
}
