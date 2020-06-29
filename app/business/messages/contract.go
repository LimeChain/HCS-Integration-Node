package messages

import "github.com/Limechain/HCS-Integration-Node/app/domain/contract/model"

type ContractMessage struct {
	BusinessMessage
	Data model.Contract `json:"data"`
}

func CreateContractRequestMessage(contract *model.Contract) *ContractMessage {
	return &ContractMessage{BusinessMessage: BusinessMessage{Type: P2PMessageTypeContractRequest}, Data: *contract}
}

func CreateContractAcceptedMessage(contract *model.Contract) *ContractMessage {
	return &ContractMessage{BusinessMessage: BusinessMessage{Type: P2PMessageTypeContractAccepted}, Data: *contract}
}
