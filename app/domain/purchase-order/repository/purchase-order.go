package repository

import "github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/model"

type PurchaseOrdersRepository interface {
	GetAll() ([]*model.PurchaseOrder, error)
	GetByID(id string) (*model.PurchaseOrder, error)
	Save(*model.PurchaseOrder) (id string, err error)
	Update(*model.PurchaseOrder) error
}
