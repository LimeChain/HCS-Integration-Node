package repository

import "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/model"

type ProposalRepository interface {
	GetAll() ([]*model.Proposal, error)
	GetByID(id string) (*model.Proposal, error)
	Save(*model.Proposal) (id string, err error)
}
