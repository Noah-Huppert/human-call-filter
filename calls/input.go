package calls

import (
	"net/http"

	"github.com/Noah-Huppert/human-call-filter/config"

	"github.com/BTBurke/twiml"
	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
)

// TestInputHandler is an http.Handler which receives digits provided by
// callers when they attempt to answer the challenge question
type TestInputHandler struct {
	// logger is used to record debug information
	logger golog.Logger

	// cfg holds application configuration
	cfg *config.Config
}

// NewTestInputHandler creates a TestInputHandler
func NewTestInputHandler(logger golog.Logger, cfg *config.Config) TestInputHandler {
	return TestInputHandler{
		logger: logger,
		cfg:    cfg,
	}
}

// ServeHTTP implements http.Handler
func (h TestInputHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get challenge answer from URL
	vars := mux.Vars(r)

	eq := vars["eq"]

	// Parse request
	var twilioReq twiml.RecordActionRequest

	err := twiml.Bind(&twilioReq, r)
	if err != nil {
		h.logger.Errorf("error reading Twilio request: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	h.logger.Debugf("received input request, correct answer: %s, request: %s",
		eq, twilioReq)

	// Get challenge response
	entered := twilioReq.Digits

	// Check response
	twilioRes := twiml.NewResponse()

	if eq == entered {
		twilioRes.Add(&twiml.Play{
			URL:    "/audio-clips/success.mp3",
			Digits: "w",
		})
		twilioRes.Add(&twiml.Dial{
			Number: h.cfg.DestinationNumber,
		})

		h.logger.Debug("caller passed challenge")
	} else {
		twilioRes.Add(&twiml.Play{
			URL:    "/audio-clips/fail.mp3",
			Digits: "w",
		})

		h.logger.Debug("caller failed challenge")
	}

	writeTwilioResp(h.logger, w, twilioRes)
}
