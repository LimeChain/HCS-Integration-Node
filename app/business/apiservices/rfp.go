package apiservices

import (
	rfpRepository "github.com/Limechain/HCS-Integration-Node/app/domain/rfp/repository"
	rfpModel "github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
)

type RFPService struct {
	repo rfpRepository.RFPRepository
}

func (s *RFPService) GetAllRFPs() ([]*rfpModel.RFP, error) {
	return s.repo.GetAll()
}

func (s *RFPService) CreateRFP(rfp *rfpModel.RFP) (id string, err error) {
	return s.repo.Save(rfp)
}

func NewRFPService(repo rfpRepository.RFPRepository) *RFPService {
	return &RFPService{repo: repo}
}