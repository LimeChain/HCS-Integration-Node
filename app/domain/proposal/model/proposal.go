package model

type Proposal struct {
	ProposalId      string       `json:"proposalId" bson:"proposalId"`
	SupplierId      string       `json:"supplierId" bson:"supplierId"`
	BuyerId         string       `json:"buyerId" bson:"buyerId"`
	ReferencedRfpId string       `json:"referencedRfpId" bson:"referencedRfpId"`
	PriceScales     []PriceScale `json:"priceScales" bson:"priceScales"`
}

type PriceScale struct {
	Sku          ProposalSku `json:"sku" bson:"sku"`
	QuantityFrom int         `json:"quantityFrom" bson:"quantityFrom"`
	QuantityTo   int         `json:"quantityTo" bson:"quantityTo"`
	SinglePrice  float32     `json:"singlePrice" bson:"singlePrice"`
	Unit         string      `json:"unit" bson:"unit"`
	Currency     string      `json:"currency" bson:"currency"`
}

type ProposalSku struct {
	ProductName       string `json:"productName" bson:"productName"`
	BuyerProductId    string `json:"buyerProductId" bson:"buyerProductId"`
	SupplierProductId string `json:"supplierProductId" bson:"supplierProductId"`
}
