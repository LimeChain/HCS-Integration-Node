package messages

type dltPO struct {
	PurchaseOrderId   string `json:"purchaseOrderId"`
	PurchaseOrderHash string `json:"purchaseOrderHash"`
	BuyerSignature    string `json:"buyerSignature"`
	SupplierSignature string `json:"supplierSignature"`
}

type DLTPOMessage struct {
	BusinessMessage
	Data dltPO `json:"data"`
}

func CreateDLTPOMessage(purchaseOrderId, purchaseOrderHash, buyerSignature, supplierSignature string) *DLTPOMessage {
	return &DLTPOMessage{
		BusinessMessage: BusinessMessage{Type: DLTMessageTypePO},
		Data: dltPO{
			PurchaseOrderId:   purchaseOrderId,
			PurchaseOrderHash: purchaseOrderHash,
			BuyerSignature:    buyerSignature,
			SupplierSignature: supplierSignature,
		},
	}
}
