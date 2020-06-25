package api

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type IntegrationNodeRoute http.Handler

type IntegrationNodeAPI struct {
	r *chi.Mux
}

type IntegrationNodeAPIResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (api *IntegrationNodeAPI) AddRouter(route string, router IntegrationNodeRoute) {
	api.r.Mount(route, router)
}

func (api *IntegrationNodeAPI) Start(port string) error {
	log.Infof("[API] Listening for REST API Requests on port %v\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), api.r)
}

func (api *IntegrationNodeAPI) StartTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, api.r)
}

func NewIntegrationNodeAPI() *IntegrationNodeAPI {
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:  []string{"*"},
		AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
	})

	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.StandardLogger(), NoColor: true}),
		middleware.Compress(6, "gzip"),
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.NoCache,
		middleware.Timeout(60*time.Second),
		c.Handler,
	)

	return &IntegrationNodeAPI{r: r}

}
