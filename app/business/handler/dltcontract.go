package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/dlt/hcs"
	log "github.com/sirupsen/logrus"
)

type DLTContractHandler struct {
	contractsRepo   repository.ContractsRepository
	contractService *service.ContractService
}

func (h *DLTContractHandler) Handle(msg *common.Message) error {

	var contractMsg messages.DLTContractMessage
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
		return errors.New("The contract was not signed by the buyer")
	}

	savedContract, err := h.contractsRepo.GetByID(contract.ContractId)
	if err != nil {
		return err
	}

	if savedContract.BuyerSignature != contract.BuyerSignature {
		return errors.New("The contract buyer signature was not the one stored")
	}

	if savedContract.SupplierSignature != contract.SupplierSignature {
		return errors.New("The contract supplier signature was not the one stored")
	}

	savedHash, err := h.contractService.Hash(&savedContract.UnsignedContract)
	if err != nil {
		return err
	}

	if savedHash != contract.ContractHash {
		return errors.New("The contract hash was not the one stored")
	}

	sn := msg.Ctx.Value(hcs.SequenceNumberKey)

	sequenceNumber, ok := sn.(uint64)
	if !ok {
		return errors.New("Could not get the proof sequence number")
	}

	savedContract.DLTAnchored = true
	savedContract.DLTProof = fmt.Sprintf("%d", sequenceNumber)

	err = h.contractsRepo.Update(savedContract)
	if err != nil {
		return err
	}

	log.Infof("Contract with Id: %s seen in the dlt and verified\n", contract.ContractId)
	return nil
}

func NewDLTContractHandler(contractsRepo repository.ContractsRepository, contractService *service.ContractService) *DLTContractHandler {
	return &DLTContractHandler{contractsRepo: contractsRepo, contractService: contractService}
}
