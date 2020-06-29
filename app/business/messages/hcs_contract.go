package messages

type blockchainContract struct {
	ContractId        string `json:"contractId"`
	ContractHash      string `json:"contractHash"`
	BuyerSignature    string `json:"buyerSignature"`
	SupplierSignature string `json:"supplierSignature"`
}

type BlockchainContractMessage struct {
	BusinessMessage
	Data blockchainContract `json:"data"`
}

func CreateBlockchainContractMessage(contractId, contractHash, buyerSignature, supplierSignature string) *BlockchainContractMessage {
	return &BlockchainContractMessage{
		BusinessMessage: BusinessMessage{Type: BlockchainMessageTypeContract},
		Data: blockchainContract{
			ContractId:        contractId,
			ContractHash:      contractHash,
			BuyerSignature:    buyerSignature,
			SupplierSignature: supplierSignature,
		},
	}
}
