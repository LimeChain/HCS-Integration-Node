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
	contractsRepo   repository.ContractsRepository
	dltClient       common.DLTMessenger
	contractService *service.ContractService
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
		return errors.New("The contract buyer signature was not the one stored. The supplier has tried to cheat you")
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

	dltMessage := messages.CreateDLTContractMessage(contract.ContractId, contractHash, contract.BuyerSignature, contract.SupplierSignature)

	dltBytes, err := json.Marshal(dltMessage)
	if err != nil {
		// TODO delete from db if cannot marshal
		return err
	}

	err = h.dltClient.Send(&common.Message{Ctx: context.TODO(), Msg: dltBytes})
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

func NewContractAcceptedHandler(contractsRepo repository.ContractsRepository, contractService *service.ContractService, dltClient common.DLTMessenger) *ContractAcceptedHandler {
	return &ContractAcceptedHandler{contractsRepo: contractsRepo, contractService: contractService, dltClient: dltClient}
}
