package dashboard

import (
	"net/http"

	"github.com/Noah-Huppert/human-call-filter/models"

	"github.com/Noah-Huppert/golog"
	"github.com/jmoiron/sqlx"
)

// ChallengesHandler is an http.Handler which returns all the challenges
// in the database
type ChallengesHandler struct {
	// logger is used to output debug information
	logger golog.Logger

	// db is a database instance
	db *sqlx.DB
}

// NewChallengesHandler creates a ChallengesHandler
func NewChallengesHandler(logger golog.Logger, db *sqlx.DB) ChallengesHandler {
	return ChallengesHandler{
		logger: logger,
		db:     db,
	}
}

// ServeHTTP implements the http.Handler
func (h ChallengesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	numbers, err := models.QueryAllChallenges(h.db)

	if err != nil {
		writeError(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(h.logger, w, http.StatusOK, map[string]interface{}{
		"challenges": numbers,
	})
}
