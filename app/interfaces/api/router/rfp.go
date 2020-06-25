package router

import (
	"github.com/Limechain/HCS-Integration-Node/app/domain/rfp/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"errors"
	log "github.com/sirupsen/logrus"
)

type storedRFPsResponse struct {
	api.IntegrationNodeAPIResponse
	RFPs []*model.RFP `json:"rfps"`
}

type createRFPResponse struct {
	api.IntegrationNodeAPIResponse
	RFPID string `json:"rfpId,omitempty"`
}

func getAllStoredRFPs(rfpService *apiservices.RFPService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storedRFPs, err := rfpService.GetAllRFPs()
		if err != nil {
			render.JSON(w, r, storedRFPsResponse{api.IntegrationNodeAPIResponse{ Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedRFPsResponse{api.IntegrationNodeAPIResponse{ Status: true, Error: "" }, storedRFPs})
	}
}

func createRFP(rfpService *apiservices.RFPService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rfp *model.RFP

		err := parser.DecodeJSONBody(w, r, &rfp)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{ Status: false, Error: mr.Msg}, ""})
				return
			}

			log.Println(err.Error())
			render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{ Status: false, Error: err.Error()}, ""})
			return
		}

		storedRFPId, errC := rfpService.CreateRFP(rfp)
		if errC != nil {
			render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{ Status: false, Error: errC.Error()}, ""})
			return
		}

		render.JSON(w, r, createRFPResponse{api.IntegrationNodeAPIResponse{ Status: true, Error: ""}, storedRFPId})
	}
}

func NewRFPRouter(rfpService *apiservices.RFPService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredRFPs(rfpService))
	r.Post("/", createRFP(rfpService))
	return r
}
