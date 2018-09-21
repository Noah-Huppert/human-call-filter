package dashboard

import (
	"net/http"

	"github.com/Noah-Huppert/human-call-filter/models"

	"github.com/Noah-Huppert/golog"
	"github.com/jmoiron/sqlx"
)

// PhoneNumbersHandler is an http.Handler which returns all the phone numbers
// in the database
type PhoneNumbersHandler struct {
	// logger is used to output debug information
	logger golog.Logger

	// db is a database instance
	db *sqlx.DB
}

// NewPhoneNumbersHandler creates a PhoneNumbersHandler
func NewPhoneNumbersHandler(logger golog.Logger, db *sqlx.DB) PhoneNumbersHandler {
	return PhoneNumbersHandler{
		logger: logger,
		db:     db,
	}
}

// ServeHTTP implements the http.Handler
func (h PhoneNumbersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	numbers, err := models.QueryAllPhoneNumbers(h.db)

	if err != nil {
		writeError(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(h.logger, w, http.StatusOK, map[string]interface{}{
		"phone_numbers": numbers,
	})
}
