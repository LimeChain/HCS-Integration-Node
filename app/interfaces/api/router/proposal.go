package router

import (
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	proposalModel "github.com/Limechain/HCS-Integration-Node/app/domain/proposal/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type CreateProposal struct {
	ProposalId      string       `json:"proposalId" bson:"proposalId"`
	SupplierId      string       `json:"supplierId" bson:"supplierId"`
	BuyerId         string       `json:"buyerId" bson:"buyerId"`
	ReferencedRfpId string       `json:"referencedRfpId" bson:"referencedRfpId"`
	PriceScales     []priceScale `json:"priceScales" bson:"priceScales"`
}

type priceScale struct {
	Sku          proposalSku `json:"sku" bson:"sku"`
	QuantityFrom int         `json:"quantityFrom" bson:"quantityFrom"`
	QuantityTo   int         `json:"quantityTo" bson:"quantityTo"`
	SinglePrice  float32     `json:"singlePrice" bson:"singlePrice"`
	Unit         string      `json:"unit" bson:"unit"`
	Currency     string      `json:"currency" bson:"currency"`
}

type proposalSku struct {
	ProductName       string `json:"productName" bson:"productName"`
	BuyerProductId    string `json:"buyerProductId" bson:"buyerProductId"`
	SupplierProductId string `json:"supplierProductId" bson:"supplierProductId"`
}

type storedProposalsResponse struct {
	api.IntegrationNodeAPIResponse
	Proposals []*proposalModel.Proposal `json:"proposals"`
}

type storedProposalResponse struct {
	api.IntegrationNodeAPIResponse
	Proposal *proposalModel.Proposal `json:"proposal"`
}

type createProposalsResponse struct {
	api.IntegrationNodeAPIResponse
	ProposalId string `json:"proposalId,omitempty"`
}

func (p *CreateProposal) toProposal() *proposalModel.Proposal {
	priceScales := make([]proposalModel.PriceScale, len(p.PriceScales))

	for i, ps := range p.PriceScales {
		sku := proposalModel.ProposalSku{
			ProductName:       ps.Sku.ProductName,
			BuyerProductId:    ps.Sku.BuyerProductId,
			SupplierProductId: ps.Sku.SupplierProductId,
		}
		priceScales[i] = proposalModel.PriceScale{
			Sku:          sku,
			QuantityFrom: ps.QuantityFrom,
			QuantityTo:   ps.QuantityTo,
			SinglePrice:  ps.SinglePrice,
			Unit:         ps.Unit,
			Currency:     ps.Currency,
		}
	}

	return &proposalModel.Proposal{
		ProposalId:      p.ProposalId,
		SupplierId:      p.SupplierId,
		BuyerId:         p.BuyerId,
		ReferencedRfpId: p.ReferencedRfpId,
		PriceScales:     priceScales,
	}
}

func getAllStoredProposals(proposalService *apiservices.ProposalService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storedProposals, err := proposalService.GetAllProposals()
		if err != nil {
			render.JSON(w, r, storedProposalsResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedProposalsResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedProposals})
	}
}

func getProposalByID(proposalService *apiservices.ProposalService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proposalID := chi.URLParam(r, "proposalID")
		storedProposal, err := proposalService.GetProposal(proposalID)
		if err != nil {
			render.JSON(w, r, storedProposalResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedProposalResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedProposal})
	}
}

func createProposal(proposalService *apiservices.ProposalService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var proposalRequest *CreateProposal

		err := parser.DecodeJSONBody(w, r, &proposalRequest)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}, ""})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		// ToDo: Validate decoded struct

		proposal := proposalRequest.toProposal()

		storedRFPId, err := proposalService.CreateProposal(proposal)
		if err != nil {
			render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedRFPId})
	}
}

func NewProposalsRouter(proposalService *apiservices.ProposalService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredProposals(proposalService))
	r.Get("/{proposalID}", getProposalByID(proposalService))
	r.Post("/", createProposal(proposalService))
	return r
}
