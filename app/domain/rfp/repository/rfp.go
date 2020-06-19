package repository

import "github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"

type RFPRepository interface {
	GetAll() ([]*model.RFP, error)
	GetByID(id string) (*model.RFP, error)
	Save(*model.RFP) error
}
