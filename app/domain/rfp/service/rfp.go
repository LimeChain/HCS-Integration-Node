package service

import (
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/repository"
)

type RFPService struct {
	r repository.RFPRepository
}

func (s *RFPService) Exists(ID string) (bool, error) {
	rfp, err := s.r.GetByID(ID)
	if err != nil {
		return false, err
	}

	return (rfp != nil), nil

}

func New(repo repository.RFPRepository) *RFPService {
	return &RFPService{r: repo}
}
