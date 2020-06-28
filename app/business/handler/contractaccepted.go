package handler

import (
	"encoding/json"
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type ContractAcceptedHandler struct {
	contractsRepo repository.ContractsRepository
	hcsClient     common.Messenger
}

func (h *ContractAcceptedHandler) Handle(msg *common.Message) error {

	var contractMsg messages.ContractMessage
	err := json.Unmarshal(msg.Msg, &contractMsg)
	if err != nil {
		return err
	}

	contract := contractMsg.Data

	// TODO add more validation

	if len(contract.BuyerSignature) == 0 {
		return errors.New("The contract was not signed by the buyer")
	}

	if len(contract.SupplierSignature) == 0 {
		return errors.New("The contract was not signed by the supplir")
	}

	// TODO Send to HCS

	err = h.contractsRepo.Update(&contract)
	if err != nil {
		return err
	}

	log.Infof("Saved accepted contract with id: %s\n", contract.ContractId)
	return nil
}

func NewContractAcceptedHandler(contractsRepo repository.ContractsRepository, hcsClient common.Messenger) *ContractAcceptedHandler {
	return &ContractAcceptedHandler{contractsRepo: contractsRepo, hcsClient: hcsClient}
}
