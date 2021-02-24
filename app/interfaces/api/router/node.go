package router

import (
	"errors"
	"net/http"

	"github.com/Limechain/HCS-Integration-Node/app/business/apiservices"
	"github.com/Limechain/HCS-Integration-Node/app/interfaces/api"
	parser "github.com/Limechain/HCS-Integration-Node/app/interfaces/api/parser"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type ConnectPeerRequest struct {
	PeerAddresses []string `json:"peerAddresses" bson:"peerAddresses"`
}

type connectPeerResponse struct {
	api.IntegrationNodeAPIResponse
}

func connectPeer(nodeService *apiservices.NodeService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var connectPeerRequest *ConnectPeerRequest

		err := parser.DecodeJSONBody(w, r, &connectPeerRequest)
		if err != nil {
			var mr *parser.MalformedRequest
			if errors.As(err, &mr) {
				log.Println(mr.Msg)
				render.JSON(w, r, connectPeerResponse{api.IntegrationNodeAPIResponse{Status: false, Error: mr.Msg}})
				return
			}

			log.Errorln(err.Error())
			render.JSON(w, r, connectPeerResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}})
			return
		}

		peerAddresses := connectPeerRequest.PeerAddresses
		for _, peerAddress := range peerAddresses {
			// ToDo: Handle the returned value
			_, err = nodeService.Connect(peerAddress)
			if err != nil {
				return
			}
		}

		if err != nil {
			render.JSON(w, r, connectPeerResponse{api.IntegrationNodeAPIResponse{Status: false, Error: err.Error()}})
			return
		}

		render.JSON(w, r, connectPeerResponse{api.IntegrationNodeAPIResponse{Status: true, Error: ""}})
	}
}

func NewNodeRouter(nodeService *apiservices.NodeService) http.Handler {
	r := chi.NewRouter()
	r.Post("/", connectPeer(nodeService))
	return r
}
