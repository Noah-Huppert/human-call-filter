package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/golog"
)

// writeJSON responds to an HTTP request with JSON
func writeJSON(logger golog.Logger, w http.ResponseWriter, status int,
	data interface{}) {

	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := encoder.Encode(data)
	if err != nil {
		logger.Errorf("error encoding JSON response: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{\"error\": \"an internal error occurred\"}")
	}
}

// writeError responds to an HTTP request with a JSON response containing an
// error field
func writeError(logger golog.Logger, w http.ResponseWriter, status int,
	err error) {

	writeJSON(logger, w, status, map[string]string{
		"error": err.Error(),
	})
}
