package service

import (
	"crypto/sha256"
	"fmt"
	"github.com/Limechain/HCS-Integration-Node/app/domain/proposal/model"
	"strings"
)

type ProposalService struct {
}

func (ps *ProposalService) Hash(p *model.Proposal) string {
	var sb strings.Builder
	sb.WriteString(p.ProposalId)
	sb.WriteRune(',')
	sb.WriteString(p.SupplierId)
	sb.WriteRune(',')
	sb.WriteString(p.BuyerId)
	sb.WriteRune(',')
	sb.WriteString(p.ReferencedRfpId)
	sb.WriteRune(',')
	for _, scale := range p.PriceScales {
		sb.WriteString(scale.Sku.ProductName)
		sb.WriteRune(',')
		sb.WriteString(scale.Sku.BuyerProductId)
		sb.WriteRune(',')
		sb.WriteString(scale.Sku.SupplierProductId)
		sb.WriteRune(',')
		sb.WriteString(fmt.Sprintf("%x", scale.QuantityFrom))
		sb.WriteRune(',')
		sb.WriteString(fmt.Sprintf("%x", scale.QuantityTo))
		sb.WriteRune(',')
		sb.WriteString(fmt.Sprintf("%x", scale.SinglePrice))
		sb.WriteRune(',')
		sb.WriteString(scale.Unit)
		sb.WriteRune(',')
		sb.WriteString(scale.Currency)
		sb.WriteRune(',')

	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(sb.String())))

}

func New() *ProposalService {
	return &ProposalService{}
}
