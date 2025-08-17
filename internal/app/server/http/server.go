package http

import (
	"github.com/aseptimu/GophKeeper/internal/app/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Server struct {
	router *chi.Mux
}

func NewServer() *Server {
	router := &Server{
		router: chi.NewRouter(),
	}
	router.registerMiddlewares()
	router.registerRoutes()
	return router
}

func (s *Server) registerMiddlewares() {
	s.router.Use(chimiddleware.Logger)
	s.router.Use(middleware.Auth)
}

func (s *Server) registerRoutes() {
	s.router.Get("/", s.handleRoot)
	s.router.Get("/auth/refresh", s.handleRoot)
}

func (s *Server) Run() {
	http.ListenAndServe(":8087", s.router)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func (s *Server) handleRefreshToken(w http.ResponseWriter, r *http.Request) {

}
