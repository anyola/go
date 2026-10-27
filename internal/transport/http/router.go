package httptransport

import (
	"net/http"

	"movieapi/internal/service"
)

func NewRouter(s *service.Service) http.Handler {
	h := NewHandler(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/load-movies", h.LoadMovies)
	mux.HandleFunc("/last-movie", h.LastMovie)
	mux.HandleFunc("/swagger/index.html", swaggerPage)
	mux.HandleFunc("/swagger/openapi.json", swaggerSpec)
	return mux
}

func swaggerPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html><html><head><title>Movie API</title><link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"></head><body><div id="swagger-ui"></div><script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script><script>window.onload=function(){SwaggerUIBundle({url:"/swagger/openapi.json",dom_id:"#swagger-ui"});}</script></body></html>`))
}
func swaggerSpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"openapi":"3.0.3","info":{"title":"Movie API","version":"1.0.0"},"paths":{"/load-movies":{"get":{"parameters":[{"name":"from","in":"query","required":true,"schema":{"type":"integer","minimum":1}},{"name":"to","in":"query","required":true,"schema":{"type":"integer","minimum":1}}],"responses":{"200":{"description":"Movies"},"400":{"description":"Invalid range"}}}},"/last-movie":{"get":{"responses":{"200":{"description":"Last movie"},"502":{"description":"Movie stand error"}}}}}}`))
}
