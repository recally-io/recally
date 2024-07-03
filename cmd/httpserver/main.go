package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"vibrain/apis/handlers"
)

func main() {
	handler := handlers.New()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8787"
	}
	s := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", port),
	}
	slog.Info("Server started", "addr", s.Addr)

	// And we serve HTTP until the world ends.
	if err := s.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
