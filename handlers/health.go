package handlers

import (
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/golog"
)

// HealthHandler provides an http.Handler for a health check endpoint
type HealthCheckHandler struct {
	// logger outputs debug information
	logger golog.Logger
}

// NewHealthHandler creates a HealthHandler
func NewHealthHandler(logger golog.Logger) HealthCheckHandler {
	return HealthCheckHandler{
		logger: logger,
	}
}

// ServeHTTP implements http.Handler
func (h HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("health check: OK")
	fmt.Fprintf(w, "OK")
}
