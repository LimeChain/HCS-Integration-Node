package model

type UnsignedPurchaseOrder struct {
	PurchaseOrderId string `json:"purchaseOrderId" bson:"purchaseOrderId"`
	SupplierId           string `json:"supplierId" bson:"supplierId"`
	BuyerId              string `json:"buyerId" bson:"buyerId"`
	ReferencedProposalId string `json:"referencedProposalId" bson:"referencedProposalId"`
	ReferencedContractId string `json:"referencedContractId" bson:"referencedContractId"`
	OrderItems []OrderItem `json:"orderItems" bson:"orderItems"`
}

type OrderItem struct {
	OrderItemId int `json:"orderItemId" bson:"orderItemId"`
	SKUBuyer string `json:"skuBuyer" bson:"skuBuyer"`
	SKUSupplier string `json:"skuSupplier" bson:"skuSupplier"`
	Quantity int `json:"quantity" bson:"quantity"`
	Unit string `json:"unit" bson:"unit"`
	SinglePrice float32 `json:"singlePrice" bson:"singlePrice"`
	TotalValue float32 `json:"totalValue" bson:"totalValue"`
	Currency string `json:"currency" bson:"currency"`
}

type PurchaseOrder struct {
	UnsignedPurchaseOrder   `json:"unsignedPurchaseOrder" bson:"unsignedPurchaseOrder"`
	BuyerSignature     string `json:"buyerSignature" bson:"buyerSignature"`
	SupplierSignature  string `json:"supplierSignature" bson:"supplierSignature"`
	BlockchainAnchored bool   `json:"blockchainAnchored" bson:"blockchainAnchored"`
	BlockchainProof    string `json:"blockchainProof" bson:"blockchainProof"`
}