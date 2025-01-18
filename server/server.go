package server

import (
	"net/http"

	"github.com/GGBusuioc/go-cache/handler"
	"github.com/GGBusuioc/go-cache/logger"
	"github.com/GGBusuioc/go-cache/router"
)

type Server struct {
	router  *router.Router
	handler *handler.Handler
	logger  logger.Logger
}

func NewServer(r *router.Router, h *handler.Handler, l logger.Logger) *Server {
	return &Server{
		router:  r,
		handler: h,
		logger:  l,
	}
}

func (s *Server) Start() error {
	mux := s.router.Setup()

	s.logger.Info("Starting server...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}
