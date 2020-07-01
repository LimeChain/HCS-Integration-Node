package handler

import (
	"encoding/json"
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type BlockchainPOHandler struct {
	por   repository.PurchaseOrdersRepository
	pos *service.PurchaseOrderService
}

func (h *BlockchainPOHandler) Handle(msg *common.Message) error {

	var purchaseOrderMsg messages.BlockchainPOMessage
	err := json.Unmarshal(msg.Msg, &purchaseOrderMsg)
	if err != nil {
		return err
	}

	po := purchaseOrderMsg.Data

	// TODO add more validation

	if len(po.BuyerSignature) == 0 {
		return errors.New("The purchase order was not signed by the buyer")
	}

	if len(po.SupplierSignature) == 0 {
		return errors.New("The purchase order was not signed by the buyer")
	}

	savedPO, err := h.por.GetByID(po.PurchaseOrderId)
	if err != nil {
		return err
	}

	if savedPO.BuyerSignature != po.BuyerSignature {
		return errors.New("The po buyer signature was not the one storred")
	}

	if savedPO.SupplierSignature != po.SupplierSignature {
		return errors.New("The po supplier signature was not the one storred")
	}

	savedHash, err := h.pos.Hash(&savedPO.UnsignedPurchaseOrder)
	if err != nil {
		return err
	}

	if savedHash != po.PurchaseOrderHash {
		return errors.New("The po hash was not the one storred")
	}

	savedPO.BlockchainAnchored = true

	err = h.por.Update(savedPO)
	if err != nil {
		return err
	}

	log.Infof("Purchase order with Id: %s seen in the blockchain and verified\n", po.PurchaseOrderId)
	return nil
}

func NewBlockchainPOHandler(por repository.PurchaseOrdersRepository, pos *service.PurchaseOrderService) *BlockchainPOHandler {
	return &BlockchainPOHandler{por: por, pos: pos}
}