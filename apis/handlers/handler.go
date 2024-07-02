package handlers

import (
	"net/http"
	"time"
	"vibrain/apis"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler is the handler for the API
type Handler struct{}

// NewServer creates a new Server with the necessary dependencies
func NewServer() apis.ServerInterface {
	options := apis.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  handleRequestError,
		ResponseErrorHandlerFunc: handleResponseError,
	}

	middlewares := make([]apis.StrictMiddlewareFunc, 0)

	return apis.NewStrictHandlerWithOptions(&Handler{}, middlewares, options)
}

func setMiddlewares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(time.Second * 60))
}

// New creates a new handler with the necessary dependencies
func New() http.Handler {
	r := chi.NewRouter()
	setMiddlewares(r)

	r.Route("/docs", func(r chi.Router) {
		r.Get("/redoc", redoc)
		r.Get("/ui", swaggerUI)
		r.Get("/json", swaggerJSON)
	})

	return apis.HandlerFromMux(NewServer(), r)
}
