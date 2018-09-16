package calls

import (
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
)

// AudioClipHandler is an http.Handler which sends audio files
type AudioClipHandler struct {
	logger golog.Logger
}

// NewAudioClipHandler creates an AudioClipHandler
func NewAudioClipHandler(logger golog.Logger) AudioClipHandler {
	return AudioClipHandler{
		logger: logger,
	}
}

// ServeHTTP implements http.Handler
func (h AudioClipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get file to serve
	vars := mux.Vars(r)

	fileName := vars["file"]

	// Serve file
	w.Header().Set("Content-Type", "audio/mpeg")
	http.ServeFile(w, r, fmt.Sprintf("./audio-clips/%s.mp3", fileName))
}
