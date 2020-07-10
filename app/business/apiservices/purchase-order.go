package apiservices

import (
	"context"
	"encoding/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	poModel "github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/model"
	poRepo "github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/repository"
	poService "github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type PurchaseOrderService struct {
	por       poRepo.PurchaseOrdersRepository
	pos       *poService.PurchaseOrderService
	p2pClient common.Messenger
}

func (s *PurchaseOrderService) GetAllPurchaseOrders() ([]*poModel.PurchaseOrder, error) {
	return s.por.GetAll()
}

func (s *PurchaseOrderService) GetPurchaseOrder(purchaseOrderId string) (*poModel.PurchaseOrder, error) {
	return s.por.GetByID(purchaseOrderId)
}

func (s *PurchaseOrderService) SaveAndSendPurchaseOrder(unsignedPurchaseOrder *poModel.UnsignedPurchaseOrder) (purchaseOrderId, purchaseOrderHash, purchaseOrderSignature string, err error) {
	purchaseOrderHash, err = s.pos.Hash(unsignedPurchaseOrder)
	if err != nil {
		return "", "", "", err
	}
	purchaseOrderSignature, err = s.pos.Sign(unsignedPurchaseOrder)
	if err != nil {
		return "", "", "", err
	}
	signedPurchaseOrder := &poModel.PurchaseOrder{UnsignedPurchaseOrder: *unsignedPurchaseOrder, BuyerSignature: purchaseOrderSignature, DLTAnchored: false}
	purchaseOrderId, err = s.por.Save(signedPurchaseOrder)
	if err != nil {
		return "", "", "", err
	}
	p2pMsg := messages.CreatePurchaseOrderRequestMessage(signedPurchaseOrder)
	p2pBytes, err := json.Marshal(p2pMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return "", "", "", err
	}
	s.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes})
	return purchaseOrderId, purchaseOrderHash, purchaseOrderSignature, nil
}

func NewPurchaseOrderService(
	por poRepo.PurchaseOrdersRepository,
	pos *poService.PurchaseOrderService,
	p2pClient common.Messenger) *PurchaseOrderService {
	return &PurchaseOrderService{por: por, pos: pos, p2pClient: p2pClient}
}
