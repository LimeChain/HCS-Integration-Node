package router

import (
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	contractModel "github.com/Limechain/HCS-Integration-Node/app/domain/contract/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type SendContractRequest struct {
	ContractId           string `json:"contractId" bson:"contractId"`
	SupplierId           string `json:"supplierId" bson:"supplierId"`
	BuyerId              string `json:"buyerId" bson:"buyerId"`
	ReferencedProposalId string `json:"referencedProposalId" bson:"referencedProposalId"`
}

type storedContractsResponse struct {
	api.IntegrationNodeAPIResponse
	Contracts []*contractModel.Contract `json:"contracts"`
}

type storedContractResponse struct {
	api.IntegrationNodeAPIResponse
	Contract *contractModel.Contract `json:"contract"`
}

type sendContractResponse struct {
	api.IntegrationNodeAPIResponse
	ContractId        string `json:"contractId, omitempty" bson:"contractId"`
	ContractHash      string `json:"contractHash, omitempty" bson:"contractHash"`
	ContractSignature string `json:"contractSignature, omitempty" bson:"contractSignature"`
}

func (req *SendContractRequest) toUnsignedContract() *contractModel.UnsignedContract {
	return &contractModel.UnsignedContract{
		ContractId:           req.ContractId,
		SupplierId:           req.SupplierId,
		BuyerId:              req.BuyerId,
		ReferencedProposalId: req.ReferencedProposalId,
	}
}

func getAllStoredContracts(contractsService *apiservices.ContractService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storedContracts, err := contractsService.GetAllContracts()
		if err != nil {
			render.JSON(w, r, storedContractsResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedContractsResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedContracts})
	}
}

func getContractById(contractsService *apiservices.ContractService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contractId := chi.URLParam(r, "contractId")
		storedContract, err := contractsService.GetContract(contractId)
		if err != nil {
			render.JSON(w, r, storedContractResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedContractResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedContract})
	}
}

func sendContract(contractsService *apiservices.ContractService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var contractRequest *SendContractRequest

		err := parser.DecodeJSONBody(w, r, &contractRequest)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, sendContractResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}, "", "", ""})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, sendContractResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, "", "", ""})
			return
		}

		// ToDo: Validate decoded struct

		unsignedContract := contractRequest.toUnsignedContract()

		contractId, contractHash, contractSignature, err := contractsService.SaveAndSendContract(unsignedContract)
		if err != nil {
			render.JSON(w, r, sendContractResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, "", "", ""})
			return
		}

		render.JSON(w, r, sendContractResponse{
			IntegrationNodeAPIResponse: api.IntegrationNodeAPIResponse{Status: true, Error: ""},
			ContractId:                 contractId,
			ContractHash:               contractHash,
			ContractSignature:          contractSignature,
		})
	}
}

func NewContractsRouter(contractsService *apiservices.ContractService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredContracts(contractsService))
	r.Get("/{contractId}", getContractById(contractsService))
	r.Post("/", sendContract(contractsService))
	return r
}
