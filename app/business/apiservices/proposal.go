package apiservices

import (
	"context"
	"encoding/json"

	"github.com/Limechain/HCS-Integration-Node/app/business/messages"
	proposalModel "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/model"
	repository "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/repository"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/common"
)

type ProposalService struct {
	repo      repository.ProposalRepository
	p2pClient common.Messenger
}

func (s *ProposalService) GetAllProposals() ([]*proposalModel.Proposal, error) {
	return s.repo.GetAll()
}

func (s *ProposalService) GetProposal(proposalId string) (*proposalModel.Proposal, error) {
	return s.repo.GetByID(proposalId)
}

func (s *ProposalService) CreateProposal(proposal *proposalModel.Proposal) (id string, err error) {
	proposalId, err := s.repo.Save(proposal)
	if err != nil {
		return "", err
	}
	p2pMsg := messages.CreateProposalMessage(proposal)
	p2pBytes, err := json.Marshal(p2pMsg)
	if err != nil {
		// TODO delete from db if cannot marshal
		return "", err
	}
	s.p2pClient.Send(&common.Message{Ctx: context.TODO(), Msg: p2pBytes}, proposal.BuyerId)
	return proposalId, nil
}

func NewProposalService(repo repository.ProposalRepository, p2pClient common.Messenger) *ProposalService {
	return &ProposalService{repo: repo, p2pClient: p2pClient}
}
