package main

import (
	"log/slog"
	"net/http"
	"vibrain/apis/handlers"
)

func main() {
	handler := handlers.New()
	s := &http.Server{
		Handler: handler,
		Addr:    ":8787",
	}
	slog.Info("Server started", "addr", s.Addr)

	// And we serve HTTP until the world ends.
	if err := s.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
