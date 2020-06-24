package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

func getAllStoredRFPs(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, struct {
		Status string `json"status"`
	}{Status: "Not Implemented"})
}

func NewRFPRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllStoredRFPs)
	return r

}
