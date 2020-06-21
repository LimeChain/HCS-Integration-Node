package repository

import "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/model"

type RFPRepository interface {
	Save(*model.Proposal) error
}
