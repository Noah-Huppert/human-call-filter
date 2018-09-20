package calls

import (
	"fmt"
	"net/http"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
)

// AudioClipsHandler is an http.Handler which sends audio files
type AudioClipsHandler struct {
	// logger prints debug information
	logger golog.Logger
}

// NewAudioClipsHandler creates an AudioClipsHandler
func NewAudioClipsHandler(logger golog.Logger) AudioClipsHandler {
	return AudioClipsHandler{
		logger: logger,
	}
}

// ServeHTTP implements http.Handler
func (h AudioClipsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get file to serve
	vars := mux.Vars(r)

	fileName := vars["file"]

	// Serve file
	w.Header().Set("Content-Type", "audio/mpeg")
	http.ServeFile(w, r, fmt.Sprintf("./audio-clips/%s.mp3", fileName))
}
