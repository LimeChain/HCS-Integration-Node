package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type PurchaseOrderRequestHandler struct {
	por       repository.PurchaseOrdersRepository
	pos       *service.PurchaseOrderService
	p2pClient common.Messenger
}

func (h *PurchaseOrderRequestHandler) Handle(msg *common.Message) error {
	remotePeerAddressCtx := msg.Ctx.Value("remotePeerAddress")

	if remotePeerAddressCtx == nil {
		return errors.New("The remote peer address is missing")
	}

	remotePeerAddress := fmt.Sprintf("%v", remotePeerAddressCtx)

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

	signatureCorrect, err := h.pos.VerifyBuyer(&po)
	if err != nil {
		return err
	}

	if !signatureCorrect {
		return errors.New("Invalid signature by the buyer")
	}

	purchaseOrderSignature, err := h.pos.Sign(&po.UnsignedPurchaseOrder)
	if err != nil {
		return err
	}

	po.SupplierSignature = purchaseOrderSignature

	purchaseOrderId, err := h.por.Save(&po)
	if err != nil {
		return err
	}

	poAcceptedMsg := messages.CreatePurchaseOrderAcceptedMessage(&po)

	p2pBytes, err := json.Marshal(poAcceptedMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return err
	}
	h.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes}, remotePeerAddress)

	log.Infof("Verified and saved po with id: %s\n", purchaseOrderId)
	return nil
}

func NewPORequestHandler(
	por repository.PurchaseOrdersRepository,
	pos *service.PurchaseOrderService,
	p2pClient common.Messenger) *PurchaseOrderRequestHandler {
	return &PurchaseOrderRequestHandler{por: por, pos: pos, p2pClient: p2pClient}
}
