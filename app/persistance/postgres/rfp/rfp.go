package rfp

import (
	"database/sql"
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
)

type RFPRepository struct {
	db *sql.DB
}

func (r *RFPRepository) GetAll() ([]*model.RFP, error) {
	return nil, nil
}

func (r *RFPRepository) GetByID(id string) (*model.RFP, error) {
	return nil, nil
}
func (r *RFPRepository) Save(*model.RFP) error {
	return nil
}

func NewRFPRepository(db *sql.DB) *RFPRepository {
	return &RFPRepository{db: db}
}
