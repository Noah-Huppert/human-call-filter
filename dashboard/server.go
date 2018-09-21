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

	// File server
	router.Handle("/api/phone_numbers", NewPhoneNumbersHandler(routeLogger, db))

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.DashboardHTTPPort),
		Handler: router,
	}
}
