package model

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type Proposal struct {
	ProposalId      string       `json:"proposalId" bson:"proposalId"`
	SupplierId      string       `json:"supplierId" bson:"supplierId"`
	BuyerId         string       `json:"buyerId" bson:"buyerId"`
	ReferencedRfpId string       `json:"referencedRfpId" bson:"referencedRfpId"`
	PriceScales     []PriceScale `json:"priceScales" bson:"priceScales"`
}

type PriceScale struct {
	Sku          ProposalSku `json:"sku" bson:"sku"`
	QuantityFrom int         `json:"quantityFrom" bson:"quantityFrom"`
	QuantityTo   int         `json:"quantityTo" bson:"quantityTo"`
	SinglePrice  float32     `json:"singlePrice" bson:"singlePrice"`
	Unit         string      `json:"unit" bson:"unit"`
	Currency     string      `json:"currency" bson:"currency"`
}

type ProposalSku struct {
	ProductName       string `json:"productName" bson:"productName"`
	BuyerProductId    string `json:"buyerProductId" bson:"buyerProductId"`
	SupplierProductId string `json:"supplierProductId" bson:"supplierProductId"`
}

func (p *Proposal) Hash() string {
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
