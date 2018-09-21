package dashboard

import (
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/golog"
	"github.com/Noah-Huppert/human-call-filter/config"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// NewServer creates a new http.Server which serves an internal API and
// dashboard page
func NewServer(logger golog.Logger, cfg *config.Config,
	db *sqlx.DB) http.Server {

	routeLogger := logger.GetChild("dashboard-routes")

	router := mux.NewRouter()

	// Phone numbers handler
	phoneNumbersLogger := routeLogger.GetChild("phone-numbers")
	phoneNumbersHandler := NewPhoneNumbersHandler(phoneNumbersLogger, db)

	router.Handle("/api/phone_numbers", phoneNumbersHandler)

	// Phone calls handler
	phoneCallsLogger := routeLogger.GetChild("phone-calls")
	phoneCallsHandler := NewPhoneCallsHandler(phoneCallsLogger, db)

	router.Handle("/api/phone_calls", phoneCallsHandler)

	// File server
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.DashboardHTTPPort),
		Handler: router,
	}
}
