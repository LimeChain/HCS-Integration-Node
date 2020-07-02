package service

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	contractRepo "github.com/Limechain/HCS-Integration-Node/app/domain/contract/repository"
	contractService "github.com/Limechain/HCS-Integration-Node/app/domain/contract/service"
	"github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/model"
	"strings"
)

type PurchaseOrderService struct {
	cr         contractRepo.ContractsRepository
	cs         *contractService.ContractService
	signingKey ed25519.PrivateKey
	peerPubKey ed25519.PublicKey
}

func (s *PurchaseOrderService) Hash(po *model.UnsignedPurchaseOrder) (string, error) {
	contract, err := s.cr.GetByID(po.ReferencedContractId)
	if err != nil {
		return "", err
	}

	contractHash, err := s.cs.Hash(&contract.UnsignedContract)

	var sb strings.Builder
	sb.WriteString(po.PurchaseOrderId)
	sb.WriteRune(',')
	sb.WriteString(po.SupplierId)
	sb.WriteRune(',')
	sb.WriteString(po.BuyerId)
	sb.WriteRune(',')
	sb.WriteString(po.ReferencedContractId)
	sb.WriteRune(',')
	sb.WriteString(contractHash)
	sb.WriteRune(',')
	sb.WriteString(contract.BuyerSignature)
	sb.WriteRune(',')
	sb.WriteString(contract.SupplierSignature)

	return fmt.Sprintf("%x", sha256.Sum256([]byte(sb.String()))), nil

}

func (s *PurchaseOrderService) Sign(po *model.UnsignedPurchaseOrder) (string, error) {
	purchaseOrderHash, err := s.Hash(po)
	if err != nil {
		return "", err
	}

	signature := ed25519.Sign(s.signingKey, []byte(purchaseOrderHash))
	signatureStr := hex.EncodeToString(signature)
	return signatureStr, nil
}

func (s *PurchaseOrderService) verify(po *model.UnsignedPurchaseOrder, signature string) (bool, error) {
	purchaseOrderHash, err := s.Hash(po)
	if err != nil {
		return false, err
	}
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}
	return ed25519.Verify(s.peerPubKey, []byte(purchaseOrderHash), signatureBytes), nil
}

func (s *PurchaseOrderService) VerifyBuyer(po *model.PurchaseOrder) (bool, error) {
	return s.verify(&po.UnsignedPurchaseOrder, po.BuyerSignature)
}

func (s *PurchaseOrderService) VerifySupplier(po *model.PurchaseOrder) (bool, error) {
	return s.verify(&po.UnsignedPurchaseOrder, po.SupplierSignature)
}

func New(signingKey ed25519.PrivateKey, cr contractRepo.ContractsRepository, cs *contractService.ContractService, peerPubKey ed25519.PublicKey) *PurchaseOrderService {
	return &PurchaseOrderService{signingKey: signingKey, cr: cr, cs: cs, peerPubKey: peerPubKey}
}
