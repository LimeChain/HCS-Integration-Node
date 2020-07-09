package router

import (
	"errors"
	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	purchaseOrderModel "github.com/Limechain/HCS-Integration-Node/app/domain/purchase-order/model"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type SendPurchaseOrderRequest struct {
	PurchaseOrderId      string             `json:"purchaseOrderId" bson:"purchaseOrderId"`
	SupplierId           string             `json:"supplierId" bson:"supplierId"`
	BuyerId              string             `json:"buyerId" bson:"buyerId"`
	ReferencedContractId string             `json:"referencedContractId" bson:"referencedContractId"`
	Items                []requestOrderItem `json:"items" bson:"items"`
}

type requestOrderItem struct {
	OrderItemId int     `json:"orderItemId" bson:"orderItemId"`
	SKUBuyer    string  `json:"skuBuyer" bson:"skuBuyer"`
	SKUSupplier string  `json:"skuSupplier" bson:"skuSupplier"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Unit        string  `json:"unit" bson:"unit"`
	SinglePrice float32 `json:"singlePrice" bson:"singlePrice"`
	TotalValue  float32 `json:"totalValue" bson:"totalValue"`
	Currency    string  `json:"currency" bson:"currency"`
}

type storedPOsResponse struct {
	api.IntegrationNodeAPIResponse
	PurchaseOrders []*purchaseOrderModel.PurchaseOrder `json:"contracts"`
}

type storedPOResponse struct {
	api.IntegrationNodeAPIResponse
	PurchaseOrder *purchaseOrderModel.PurchaseOrder `json:"po"`
}

type sendPOResponse struct {
	api.IntegrationNodeAPIResponse
	PurchaseOrderId        string `json:"purchaseOrderId, omitempty" bson:"purchaseOrderId"`
	PurchaseOrderHash      string `json:"purchaseOrderHash, omitempty" bson:"purchaseOrderHash"`
	PurchaseOrderSignature string `json:"purchaseOrderSignature, omitempty" bson:"purchaseOrderSignature"`
}

func (req *SendPurchaseOrderRequest) toUnsignedPurchaseOrder() *purchaseOrderModel.UnsignedPurchaseOrder {
	items := make([]purchaseOrderModel.OrderItem, len(req.Items))

	for i, item := range req.Items {
		items[i] = purchaseOrderModel.OrderItem{
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

	return &purchaseOrderModel.UnsignedPurchaseOrder{
		PurchaseOrderId:      req.PurchaseOrderId,
		SupplierId:           req.SupplierId,
		BuyerId:              req.BuyerId,
		ReferencedContractId: req.ReferencedContractId,
		OrderItems:           items,
	}
}

func getAllStoredPOs(poService *apiservices.PurchaseOrderService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		storedPOs, err := poService.GetAllPurchaseOrders()
		if err != nil {
			render.JSON(w, r, storedPOsResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedPOsResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedPOs})
	}
}

func getPurchaseOrderById(poService *apiservices.PurchaseOrderService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		purchaseOrderId := chi.URLParam(r, "purchaseOrderId")
		storedPO, err := poService.GetPurchaseOrder(purchaseOrderId)
		if err != nil {
			render.JSON(w, r, storedPOResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, nil})
			return
		}
		render.JSON(w, r, storedPOResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}, storedPO})
	}
}

func sendPO(poService *apiservices.PurchaseOrderService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var poRequest *SendPurchaseOrderRequest

		err := parser.DecodeJSONBody(w, r, &poRequest)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, sendPOResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}, "", "", ""})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, sendPOResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, "", "", ""})
			return
		}

		// ToDo: Validate decoded struct

		unsignedPO := poRequest.toUnsignedPurchaseOrder()

		purchaseOrderId, purchaseOrderHash, purchaseOrderSignature, err := poService.SaveAndSendPurchaseOrder(unsignedPO)
		if err != nil {
			render.JSON(w, r, sendPOResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}, "", "", ""})
			return
		}

		render.JSON(w, r, sendPOResponse{
			IntegrationNodeAPIResponse: api.IntegrationNodeAPIResponse{Status: true, Error: ""},
			PurchaseOrderId:            purchaseOrderId,
			PurchaseOrderHash:          purchaseOrderHash,
			PurchaseOrderSignature:     purchaseOrderSignature,
		})
	}
}

func NewPurchaseOrdersRouter(purchaseOrdersService *apiservices.PurchaseOrderService) http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredPOs(purchaseOrdersService))
	r.Get("/{purchaseOrderId}", getPurchaseOrderById(purchaseOrdersService))
	r.Post("/", sendPO(purchaseOrdersService))
	return r
}
