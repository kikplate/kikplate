package lib

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RequestHandler struct {
	Mux *chi.Mux
}

func NewRequestHandler(logger Logger) RequestHandler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)

	return RequestHandler{Mux: mux}
}
