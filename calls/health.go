package calls

import (
	"fmt"
	"net/http"
)

// HealthHandler provides an http.Handler for a health check endpoint
type HealthCheckHandler struct{}

// NewHealthHandler creates a HealthHandler
func NewHealthHandler() HealthCheckHandler {
	return HealthCheckHandler{}
}

// ServeHTTP implements http.Handler
func (h HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
