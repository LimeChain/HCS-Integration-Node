package apiservices

import (
	"context"
	"encoding/json"
	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	rfpModel "github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
	rfpRepository "github.com/Limechain/HCS-Integration-Node/app/domain/rfp/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type RFPService struct {
	repo      rfpRepository.RFPRepository
	p2pClient common.Messenger
}

func (s *RFPService) GetAllRFPs() ([]*rfpModel.RFP, error) {
	return s.repo.GetAll()
}

func (s *RFPService) CreateRFP(rfp *rfpModel.RFP) (id string, err error) {
	rfpId, err := s.repo.Save(rfp)
	if err != nil {
		return "", err
	}
	p2pMsg := messages.CreateRFPMessage(rfp)
	p2pBytes, err := json.Marshal(p2pMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return "", err
	}
	s.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes})
	return rfpId, nil
}

func NewRFPService(repo rfpRepository.RFPRepository, p2pClient common.Messenger) *RFPService {
	return &RFPService{repo: repo, p2pClient: p2pClient}
}
