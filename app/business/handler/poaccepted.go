package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type PurchaseOrderAcceptedHandler struct {
	por  repository.PurchaseOrdersRepository
	blockchainClient common.Messenger
	pos  *service.PurchaseOrderService
}

func (h *PurchaseOrderAcceptedHandler) Handle(msg *common.Message) error {

	var purchaseOrderMsg messages.PurchaseOrderMessage
	err := json.Unmarshal(msg.Msg, &purchaseOrderMsg)
	if err != nil {
		return err
	}

	po := purchaseOrderMsg.Data

	// TODO add more validation

	if len(po.BuyerSignature) == 0 {
		return errors.New("The po was not signed by the buyer")
	}

	if len(po.SupplierSignature) == 0 {
		return errors.New("The po was not signed by the supplir")
	}

	savedPO, err := h.por.GetByID(po.PurchaseOrderId)
	if err != nil {
		return err
	}

	if savedPO.BuyerSignature != po.BuyerSignature {
		return errors.New("The po buyer signature was not the one storred. The supplier has tried to cheat you")
	}

	purchaseOrderHash, err := h.pos.Hash(&po.UnsignedPurchaseOrder)
	if err != nil {
		return err
	}

	signatureCorrect, err := h.pos.VerifySupplier(&po)
	if err != nil {
		return err
	}

	if !signatureCorrect {
		return errors.New("Invalid signature by the supplier")
	}

	blockchainMessage := messages.CreateBlockchainPOMessage(po.PurchaseOrderId, purchaseOrderHash, po.BuyerSignature, po.SupplierSignature)

	blockchainBytes, err := json.Marshal(blockchainMessage)
	if err != nil {
		// TODO delete from db if cannot marshal
		return err
	}

	err = h.blockchainClient.Send(&common.Message{Ctx: context.TODO(), Msg: blockchainBytes})
	if err != nil {
		return err
	}

	err = h.por.Update(&po)
	if err != nil {
		return err
	}

	log.Infof("Verified and saved accepted po with id: %s\n", po.PurchaseOrderId)
	return nil
}

func NewPOAcceptedHandler(por repository.PurchaseOrdersRepository, pos *service.PurchaseOrderService, blockchainClient common.Messenger) *PurchaseOrderAcceptedHandler {
	return &PurchaseOrderAcceptedHandler{por: por, pos: pos, blockchainClient: blockchainClient}
}