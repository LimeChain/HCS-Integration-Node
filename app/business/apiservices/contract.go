package apiservices

import (
	"context"
	"encoding/json"

	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	contractModel "github.com/Limechain/HCS-Integration-Node/app/domain/contract/model"
	contractRepo "github.com/Limechain/HCS-Integration-Node/app/domain/contract/repository"
	contractService "github.com/Limechain/HCS-Integration-Node/app/domain/contract/service"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type ContractService struct {
	cr        contractRepo.ContractsRepository
	cs        *contractService.ContractService
	p2pClient common.Messenger
}

func (s *ContractService) GetAllContracts() ([]*contractModel.Contract, error) {
	return s.cr.GetAll()
}

func (s *ContractService) GetContract(contractId string) (*contractModel.Contract, error) {
	return s.cr.GetByID(contractId)
}

func (s *ContractService) SaveAndSendContract(unsignedContract *contractModel.UnsignedContract) (contractId, contractHash, contractSignature string, err error) {
	contractHash, err = s.cs.Hash(unsignedContract)
	if err != nil {
		return "", "", "", err
	}
	contractSignature, err = s.cs.Sign(unsignedContract)
	if err != nil {
		return "", "", "", err
	}
	signedContract := &contractModel.Contract{UnsignedContract: *unsignedContract, BuyerSignature: contractSignature, DLTAnchored: false}
	contractId, err = s.cr.Save(signedContract)
	if err != nil {
		return "", "", "", err
	}
	p2pMsg := messages.CreateContractRequestMessage(signedContract)
	p2pBytes, err := json.Marshal(p2pMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return "", "", "", err
	}
	s.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes}, signedContract.SupplierId)
	return contractId, contractHash, contractSignature, nil
}

func NewContractService(
	cr contractRepo.ContractsRepository,
	cs *contractService.ContractService,
	p2pClient common.Messenger) *ContractService {
	return &ContractService{cr: cr, cs: cs, p2pClient: p2pClient}
}
