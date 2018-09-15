package calls

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/BTBurke/twiml"
	"github.com/Noah-Huppert/golog"
)

// CallsHandler implements a http.Handler which responds to Twilio phone calls
type CallsHandler struct {
	logger golog.Logger
}

// NewCallsHandler creates a new CallsHandler
func NewCallsHandler(logger golog.Logger) CallsHandler {
	return CallsHandler{
		logger: logger,
	}
}

// writeXML writes a Twilio response as an HTTP response
func writeTwilioResp(logger golog.Logger, w http.ResponseWriter,
	twilioRes *twiml.Response) {

	// Encode Twilio response to bytes
	bytes, err := twilioRes.Encode()

	if err != nil {
		logger.Errorf("error encoding twilio response into bytes: %s",
			err.Error())

		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// Write bytes as response
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bytes)
	if err != nil {
		logger.Errorf("error writing twilio response: %s", err.Error())

		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

// ServeHTTP handles a Twilio phone call
func (h CallsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var twilioReq twiml.VoiceRequest

	err := twiml.Bind(&twilioReq, r)
	if err != nil {
		h.logger.Errorf("error reading Twilio request: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	// Create response
	twilioRes := twiml.NewResponse()

	switch twilioReq.CallStatus {
	case twiml.InProgress:
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, http.StatusText(http.StatusOK))

	case twiml.Ringing, twiml.Queued:
		a := rand.Intn(5)
		b := rand.Intn(5)
		//eq := a + b

		twilioRes.Add(&twiml.Play{
			URL:    "http://f062d053.ngrok.io/audio-clips/intro.mp3",
			Digits: "w",
		})
		twilioRes.Add(&twiml.Say{
			Voice: "woman",
			Text:  fmt.Sprintf("What is %d plus %d?", a, b),
		})
		writeTwilioResp(h.logger, w, twilioRes)

	default:
		twilioRes.Add(&twiml.Hangup{})
		writeTwilioResp(h.logger, w, twilioRes)
	}
}
