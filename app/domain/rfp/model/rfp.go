package model

type RFP struct {
	RFPId  string `json:"rfpId" bson:"rfpId"`
	SupplierId string `json:"supplierId" bson:"supplierId"`
	BuyerId string `json:"buyerId" bson:"buyerId"`
	Items []Item `json:"items" bson:"items"`
}

type Item struct {
	OrderItemId int `json:"orderItemId" bson:"orderItemId"`
	SKUBuyer string `json:"skuBuyer" bson:"skuBuyer"`
	SKUSupplier string `json:"skuSupplier" bson:"skuSupplier"`
	Quantity int `json:"quantity" bson:"quantity"`
	Unit string `json:"unit" bson:"unit"`
	SinglePrice float32 `json:"singlePrice" bson:"singlePrice"`
	TotalValue float32 `json:"totalValue" bson:"totalValue"`
	Currency string `json:"currency" bson:"currency"`
}