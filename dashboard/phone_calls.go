package dashboard

import (
	"net/http"

	"github.com/Noah-Huppert/human-call-filter/models"

	"github.com/Noah-Huppert/golog"
	"github.com/jmoiron/sqlx"
)

// PhoneCallsHandler is an http.Handler which returns all the phone calls in
// the database
type PhoneCallsHandler struct {
	// logger is used to output debug information
	logger golog.Logger

	// db is a database instance
	db *sqlx.DB
}

// NewPhoneCallsHandler creates a PhoneCallsHandler
func NewPhoneCallsHandler(logger golog.Logger, db *sqlx.DB) PhoneCallsHandler {
	return PhoneCallsHandler{
		logger: logger,
		db:     db,
	}
}

// ServeHTTP implements the http.Handler
func (h PhoneCallsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	calls, err := models.QueryAllPhoneCalls(h.db)

	if err != nil {
		writeError(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(h.logger, w, http.StatusOK, map[string]interface{}{
		"phone_calls": calls,
	})
}
