package http

import (
	"github.com/aseptimu/GophKeeper/internal/app/config"
	"github.com/aseptimu/GophKeeper/internal/app/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

type Registrer interface {
	RegisterRoutes(r chi.Router)
}

type Server struct {
	router *chi.Mux
	config *config.AppConfig
}

func NewServer(config *config.AppConfig, routes ...Registrer) *Server {
	server := &Server{
		router: chi.NewRouter(),
		config: config,
	}
	server.registerMiddlewares()
	for _, reg := range routes {
		reg.RegisterRoutes(server.router)
	}
	return server
}

func (s *Server) registerMiddlewares() {
	s.router.Use(chimiddleware.Logger)
	s.router.Use(middleware.Auth([]byte(s.config.JWTKey)))
}

func (s *Server) Run() {
	err := http.ListenAndServe(s.config.ServerAddress, s.router)
	if err != nil {
		slog.Error("Server failure", "error", err)
	}
}
