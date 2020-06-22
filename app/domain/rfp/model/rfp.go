package model

type RFP struct {
	RFPID  string `json:"rfpId" bson:"rfpId"`
	ItemId string `json:"itemId" bson:"itemId,omitempty"`
}
