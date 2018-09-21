package dashboard

import (
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/human-call-filter/config"
	//"github.com/gorilla/mux"
)

// NewServer creates a new http.Server which serves an internal API and
// dashboard page
func NewServer(cfg *config.Config) http.Server {
	//router := mux.NewRouter()

	// File server
	//router.Handle("/", http.FileServer(http.Dir("static")))

	return http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.DashboardHTTPPort),
		Handler: http.FileServer(http.Dir("static")),
	}
}
