package server

import (
	"log"
	"net/http"

	cfg "github.com/MrR0b0t1001/Chirpy/config"
	"github.com/MrR0b0t1001/Chirpy/handlers/health"
	"github.com/MrR0b0t1001/Chirpy/utils"
)

type APIServer struct {
	listenAddr string
	handler    *http.ServeMux
	config     *cfg.APIConfig
}

func NewAPIServer(listenAddr string, handler *http.ServeMux, config *cfg.APIConfig) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		handler:    handler,
		config:     config,
	}
}

func (s *APIServer) Run() {
	s.handler.Handle(
		"GET /app/",
		s.config.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))),
	)

	s.handler.Handle(
		"GET /app/assets/",
		s.config.MiddlewareMetricsInc(
			http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets"))),
		),
	)

	s.handler.HandleFunc("GET /api/healthz", health.HandleHealthz)
	s.handler.HandleFunc("GET /api/chirps", utils.MakeHTTPHandleFunc(s.config.HandleGetChirps))
	s.handler.HandleFunc(
		"GET /api/chirps/{chirpID}",
		utils.MakeHTTPHandleFunc(s.config.HandleGetChirpByID),
	)

	s.handler.HandleFunc("GET /admin/metrics", utils.MakeHTTPHandleFunc(s.config.MetricsHandler))
	s.handler.HandleFunc("POST /admin/reset", utils.MakeHTTPHandleFunc(s.config.HandleDeleteUsers))

	s.handler.HandleFunc("POST /api/users", utils.MakeHTTPHandleFunc(s.config.HandleCreateUser))
	s.handler.HandleFunc("POST /api/chirps", utils.MakeHTTPHandleFunc(s.config.HandleCreateChirp))
	s.handler.HandleFunc("POST /api/login", utils.MakeHTTPHandleFunc(s.config.HandleLogin))
	s.handler.HandleFunc("POST /api/refresh", utils.MakeHTTPHandleFunc(s.config.HandleRefresh))
	s.handler.HandleFunc("POST /api/revoke", utils.MakeHTTPHandleFunc(s.config.HandleRevoke))

	s.handler.HandleFunc("PUT /api/users", utils.MakeHTTPHandleFunc(s.config.HandleUpdateUser))

	log.Printf("Starting server on %s...", s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, s.handler); err != nil {
		log.Printf("Error starting the server: %v\n", err)
	}
}
