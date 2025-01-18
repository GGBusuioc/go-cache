package router

import (
	"net/http"

	"github.com/GGBusuioc/go-cache/handler"
	"github.com/GGBusuioc/go-cache/logger"
)

type Router struct {
	handler *handler.Handler
	logger  logger.Logger
}

func NewRouter(h *handler.Handler, l logger.Logger) *Router {
	return &Router{
		handler: h,
		logger:  l,
	}
}

func (r *Router) Setup() *http.ServeMux {
	r.logger.Debug("Setting up routes...")

	// create a new request multiplexer
	// take incoming requests and display them to the matching handlers
	mux := http.NewServeMux()

	// register the routes and handlers
	mux.Handle("/", &handler.HomeHandler{})
	mux.Handle("/cache", r.handler)
	mux.Handle("/cache/", r.handler)

	return mux
}
