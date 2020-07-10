package model

type UnsignedContract struct {
	ContractId           string `json:"contractId" bson:"contractId"`
	SupplierId           string `json:"supplierId" bson:"supplierId"`
	BuyerId              string `json:"buyerId" bson:"buyerId"`
	ReferencedProposalId string `json:"referencedProposalId" bson:"referencedProposalId"`
}

type Contract struct {
	UnsignedContract  `json:"unsignedContract" bson:"unsignedContract"`
	BuyerSignature    string `json:"buyerSignature" bson:"buyerSignature"`
	SupplierSignature string `json:"supplierSignature" bson:"supplierSignature"`
	DLTAnchored       bool   `json:"DLTAnchored" bson:"DLTAnchored"`
	DLTProof          string `json:"DLTProof" bson:"DLTProof"`
}
