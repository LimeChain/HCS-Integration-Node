package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type ContractAcceptedHandler struct {
	contractsRepo    repository.ContractsRepository
	blockchainClient common.Messenger
	contractService  *service.ContractService
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

	savedContract, err := h.contractsRepo.GetByID(contract.ContractId)
	if err != nil {
		return err
	}

	if savedContract.BuyerSignature != contract.BuyerSignature {
		return errors.New("The contract buyer signature was not the one storred. The supplier has tried to cheat you")
	}

	contractHash, err := h.contractService.Hash(&contract.UnsignedContract)
	if err != nil {
		return err
	}

	signatureCorrect, err := h.contractService.VerifySupplier(&contract)
	if err != nil {
		return err
	}

	if !signatureCorrect {
		return errors.New("Invalid signature by the supplier")
	}

	blockchainMessage := messages.CreateBlockchainContractMessage(contract.ContractId, contractHash, contract.BuyerSignature, contract.SupplierSignature)

	blockchainBytes, err := json.Marshal(blockchainMessage)
	if err != nil {
		// TODO delete from db if cannot marshal
		return err
	}

	err = h.blockchainClient.Send(&common.Message{Ctx: context.TODO(), Msg: blockchainBytes})
	if err != nil {
		return err
	}

	err = h.contractsRepo.Update(&contract)
	if err != nil {
		return err
	}

	log.Infof("Verified and saved accepted contract with id: %s\n", contract.ContractId)
	return nil
}

func NewContractAcceptedHandler(contractsRepo repository.ContractsRepository, contractService *service.ContractService, blockchainClient common.Messenger) *ContractAcceptedHandler {
	return &ContractAcceptedHandler{contractsRepo: contractsRepo, contractService: contractService, blockchainClient: blockchainClient}
}
