package router

import (
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type CreateRFPRequest struct {
	RFPId      string        `json:"rfpId" bson:"rfpId"`
	SupplierId string        `json:"supplierId" bson:"supplierId"`
	BuyerId    string        `json:"buyerId" bson:"buyerId"`
	Items      []requestItem `json:"items" bson:"items"`
}

type requestItem struct {
	OrderItemId int     `json:"orderItemId" bson:"orderItemId"`
	SKUBuyer    string  `json:"skuBuyer" bson:"skuBuyer"`
	SKUSupplier string  `json:"skuSupplier" bson:"skuSupplier"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Unit        string  `json:"unit" bson:"unit"`
	SinglePrice float32 `json:"singlePrice" bson:"singlePrice"`
	TotalValue  float32 `json:"totalValue" bson:"totalValue"`
	Currency    string  `json:"currency" bson:"currency"`
}

type storedRFPsResponse struct {
	api.IntegrationNodeAPIResponse
	RFPs []*model.RFP `json:"rfps"`
}

type storedRFPResponse struct {
	api.IntegrationNodeAPIResponse
	RFP *model.RFP `json:"rfp"`
}

type createRFPResponse struct {
	api.IntegrationNodeAPIResponse
	RFPId string `json:"rfpId,omitempty"`
}

func (rfpRequestModel *CreateRFPRequest) toRFP() *model.RFP {
	items := make([]model.Item, len(rfpRequestModel.Items))

	for i, item := range rfpRequestModel.Items {
		items[i] = model.Item{
			OrderItemId: item.OrderItemId,
			SKUBuyer:    item.SKUBuyer,
			SKUSupplier: item.SKUSupplier,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			SinglePrice: item.SinglePrice,
			TotalValue:  item.TotalValue,
			Currency:    item.Currency,
		}
	}

	return &model.RFP{
		RFPId:      rfpRequestModel.RFPId,
		SupplierId: rfpRequestModel.SupplierId,
		BuyerId:    rfpRequestModel.BuyerId,
		Items:      items,
	}
}

func getAllStoredRFPs(rfpService *apiservices.RFPService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storedRFPs, err := rfpService.GetAllRFPs()
		if err != nil {
			render.JSON(w, r, storedRFPsResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedRFPsResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedRFPs})
	}
}

func getRFPById(rfpService *apiservices.RFPService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rfpId := chi.URLParam(r, "rfpId")
		rfp, err := rfpService.GetRFP(rfpId)
		if err != nil {
			render.JSON(w, r, storedRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedRFPResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, rfp})
	}
}

func createRFP(rfpService *apiservices.RFPService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rfpRequest *CreateRFPRequest

		err := parser.DecodeJSONBody(w, r, &rfpRequest)
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

		rfp := rfpRequest.toRFP()

		storedRFPId, err := rfpService.CreateRFP(rfp)
		if err != nil {
			render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, ""})
			return
		}

		render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedRFPId})
	}
}

func NewRFPRouter(rfpService *apiservices.RFPService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredRFPs(rfpService))
	r.Get("/{rfpId}", getRFPById(rfpService))
	r.Post("/", createRFP(rfpService))
	return r
}
