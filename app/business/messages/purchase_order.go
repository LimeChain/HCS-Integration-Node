package messages

import "github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/model"

type PurchaseOrderMessage struct {
	BusinessMessage
	Data model.PurchaseOrder `json:"data"`
}

func CreatePurchaseOrderRequestMessage(po *model.PurchaseOrder) *PurchaseOrderMessage {
	return &PurchaseOrderMessage{BusinessMessage: BusinessMessage{Type: P2PMessageTypePORequest}, Data: *po}
}

func CreatePurchaseOrderAcceptedMessage(po *model.PurchaseOrder) *PurchaseOrderMessage {
	return &PurchaseOrderMessage{BusinessMessage: BusinessMessage{Type: P2PMessageTypePOAccepted}, Data: *po}
}
