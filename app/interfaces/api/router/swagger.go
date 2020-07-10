package router

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func NewSwaggerRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("./docs/swagger.json"), //The url pointing to API definition"
	))

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "swagger"))
	FileServer(r, "/docs", filesDir)

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}