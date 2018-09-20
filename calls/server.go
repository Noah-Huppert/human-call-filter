package calls

import (
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/human-call-filter/config"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// NewServer creates a http.Server which handles twilio call requests
func NewServer(logger golog.Logger, cfg *config.Config,
	db *sqlx.DB) http.Server {

	router := mux.NewRouter()

	routesLogger := logger.GetChild("call-routes")

	// Health check handler
	healthHandler := NewHealthHandler()

	router.Handle("/healthz", healthHandler)

	// Call handler
	callLogger := routesLogger.GetChild("call")
	callHandler := NewCallsHandler(callLogger, cfg, db)

	router.Handle("/call", callHandler).Methods("POST")

	// Challenge input handler
	inputLogger := routesLogger.GetChild("input")
	inputHandler := NewTestInputHandler(inputLogger, cfg, db)

	router.Handle("/input/challenge/{challenge_id}", inputHandler).
		Methods("POST")

	// Audio clips handler
	audioClipLogger := routesLogger.GetChild("audio-clips")
	audioClipHandler := NewAudioClipsHandler(audioClipLogger)

	router.Handle("/audio-clips/{file}.mp3", audioClipHandler)

	return http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: router,
	}
}
