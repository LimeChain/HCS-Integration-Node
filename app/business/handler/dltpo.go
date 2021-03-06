package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/dlt/hcs"
	log "github.com/sirupsen/logrus"
)

type DLTPOHandler struct {
	por repository.PurchaseOrdersRepository
	pos *service.PurchaseOrderService
}

func (h *DLTPOHandler) Handle(msg *common.Message) error {

	var purchaseOrderMsg messages.DLTPOMessage
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
		return errors.New("The po buyer signature was not the one stored")
	}

	if savedPO.SupplierSignature != po.SupplierSignature {
		return errors.New("The po supplier signature was not the one stored")
	}

	savedHash, err := h.pos.Hash(&savedPO.UnsignedPurchaseOrder)
	if err != nil {
		return err
	}

	sn := msg.Ctx.Value(hcs.SequenceNumberKey)

	sequenceNumber, ok := sn.(uint64)
	if !ok {
		return errors.New("Could not get the proof sequence number")
	}

	if savedHash != po.PurchaseOrderHash {
		return errors.New("The po hash was not the one stored")
	}

	savedPO.DLTAnchored = true
	savedPO.DLTProof = fmt.Sprintf("%d", sequenceNumber)

	err = h.por.Update(savedPO)
	if err != nil {
		return err
	}

	log.Infof("Purchase order with Id: %s seen in the dlt and verified\n", po.PurchaseOrderId)
	return nil
}

func NewDLTPOHandler(por repository.PurchaseOrdersRepository, pos *service.PurchaseOrderService) *DLTPOHandler {
	return &DLTPOHandler{por: por, pos: pos}
}
