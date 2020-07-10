package messages

type dltContract struct {
	ContractId        string `json:"contractId"`
	ContractHash      string `json:"contractHash"`
	BuyerSignature    string `json:"buyerSignature"`
	SupplierSignature string `json:"supplierSignature"`
}

type DLTContractMessage struct {
	BusinessMessage
	Data dltContract `json:"data"`
}

func CreateDLTContractMessage(contractId, contractHash, buyerSignature, supplierSignature string) *DLTContractMessage {
	return &DLTContractMessage{
		BusinessMessage: BusinessMessage{Type: DLTMessageTypeContract},
		Data: dltContract{
			ContractId:        contractId,
			ContractHash:      contractHash,
			BuyerSignature:    buyerSignature,
			SupplierSignature: supplierSignature,
		},
	}
}
