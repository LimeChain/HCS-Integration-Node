package messages

type blockchainPO struct {
	PurchaseOrderId   string `json:"purchaseOrderId"`
	PurchaseOrderHash string `json:"purchaseOrderHash"`
	BuyerSignature    string `json:"buyerSignature"`
	SupplierSignature string `json:"supplierSignature"`
}

type BlockchainPOMessage struct {
	BusinessMessage
	Data blockchainPO `json:"data"`
}

func CreateBlockchainPOMessage(purchaseOrderId, purchaseOrderHash, buyerSignature, supplierSignature string) *BlockchainPOMessage {
	return &BlockchainPOMessage{
		BusinessMessage: BusinessMessage{Type: BlockchainMessageTypePO},
		Data: blockchainPO{
			PurchaseOrderId:        purchaseOrderId,
			PurchaseOrderHash:      purchaseOrderHash,
			BuyerSignature:    buyerSignature,
			SupplierSignature: supplierSignature,
		},
	}
}
