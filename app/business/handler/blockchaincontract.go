package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/repository"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/blockchain/hcs"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
	log "github.com/sirupsen/logrus"
)

type BlockchainContractHandler struct {
	contractsRepo   repository.ContractsRepository
	contractService *service.ContractService
}

func (h *BlockchainContractHandler) Handle(msg *common.Message) error {

	var contractMsg messages.BlockchainContractMessage
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
		return errors.New("The contract buyer signature was not the one storred")
	}

	if savedContract.SupplierSignature != contract.SupplierSignature {
		return errors.New("The contract supplier signature was not the one storred")
	}

	savedHash, err := h.contractService.Hash(&savedContract.UnsignedContract)
	if err != nil {
		return err
	}

	if savedHash != contract.ContractHash {
		return errors.New("The contract hash was not the one storred")
	}

	sn := msg.Ctx.Value(hcs.SequenceNumberKey)

	sequenceNumber, ok := sn.(uint64)
	if !ok {
		return errors.New("Could not get the proof sequence number")
	}

	savedContract.BlockchainAnchored = true
	savedContract.BlockchainProof = fmt.Sprintf("%d", sequenceNumber)

	err = h.contractsRepo.Update(savedContract)
	if err != nil {
		return err
	}

	log.Infof("Contract with Id: %s seen in the blockchain and verified\n", contract.ContractId)
	return nil
}

func NewBlockchainContractHandler(contractsRepo repository.ContractsRepository, contractService *service.ContractService) *BlockchainContractHandler {
	return &BlockchainContractHandler{contractsRepo: contractsRepo, contractService: contractService}
}
