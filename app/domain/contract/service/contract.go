package service

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/domain/contract/model"
	proposalRepo "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/repository"
	proposalService "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/service"
	"strings"
)

type ContractService struct {
	pr         proposalRepo.ProposalRepository
	ps         proposalService.ProposalService
	signingKey ed25519.PrivateKey
}

func (s *ContractService) Hash(c *model.UnsignedContract) (string, error) {
	proposal, err := s.pr.GetByID(c.ReferencedProposalId)
	if err != nil {
		return "", err
	}

	proposalHash := s.ps.Hash(proposal)

	var sb strings.Builder
	sb.WriteString(c.ContractId)
	sb.WriteRune(',')
	sb.WriteString(c.SupplierId)
	sb.WriteRune(',')
	sb.WriteString(c.BuyerId)
	sb.WriteRune(',')
	sb.WriteString(c.ReferencedProposalId)
	sb.WriteRune(',')
	sb.WriteString(proposalHash)

	return fmt.Sprintf("%x", sha256.Sum256([]byte(sb.String()))), nil

}

func (s *ContractService) Sign(c *model.UnsignedContract) (string, error) {
	contractHash, err := s.Hash(c)
	if err != nil {
		return "", err
	}

	signature := ed25519.Sign(s.signingKey, []byte(contractHash))
	return string(signature), nil
}

func New(signingKey ed25519.PrivateKey, pr proposalRepo.ProposalRepository, ps proposalService.ProposalService) *ContractService {
	return &ContractService{signingKey: signingKey, pr: pr, ps: ps}
}
