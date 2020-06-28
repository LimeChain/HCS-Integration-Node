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

type ContractRequestHandler struct {
	contractsRepo    repository.ContractsRepository
	contractsService *service.ContractService
	p2pClient        common.Messenger
}

func (h *ContractRequestHandler) Handle(msg *common.Message) error {

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

	contractSignature, err := h.contractsService.Sign(&contract.UnsignedContract)
	if err != nil {
		return err
	}

	contract.SupplierSignature = contractSignature

	contractId, err := h.contractsRepo.Save(&contract)
	if err != nil {
		return err
	}

	contractAcceptedMsg := messages.CreateContractAcceptedMessage(&contract)

	p2pBytes, err := json.Marshal(contractAcceptedMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return err
	}
	h.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes})

	log.Infof("Saved contract with id: %s\n", contractId)
	return nil
}

func NewContractRequestHandler(contractsRepo repository.ContractsRepository, contractsService *service.ContractService, p2pClient common.Messenger) *ContractRequestHandler {
	return &ContractRequestHandler{contractsRepo: contractsRepo, contractsService: contractsService, p2pClient: p2pClient}
}
