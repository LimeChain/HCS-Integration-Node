package repository

import "github.com/Limechain/HCS-Integration-Node/app/domain/contract/model"

type ContractsRepository interface {
	GetAll() ([]*model.Contract, error)
	GetByID(id string) (*model.Contract, error)
	Save(*model.Contract) (id string, err error)
	Update(*model.Contract) error
}
