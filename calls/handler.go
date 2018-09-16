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

	h.logger.Debugf("received call request %#v", twilioReq)

	switch twilioReq.CallStatus {
	case twiml.InProgress:
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, http.StatusText(http.StatusOK))

	case twiml.Ringing, twiml.Queued:
		a := rand.Intn(3) + 1
		b := rand.Intn(3) + 1
		eq := a + b

		twilioRes.Add(&twiml.Play{
			URL:    "/audio-clips/intro.mp3",
			Digits: "w",
		})
		twilioRes.Add(&twiml.Say{
			Voice: "man",
			Text:  fmt.Sprintf("What is ,%d. plus ,%d", a, b),
		})
		twilioRes.Add(&twiml.Gather{
			Action:    fmt.Sprintf("/input/test/%d", eq),
			NumDigits: 1,
			Timeout:   10,
		})
		writeTwilioResp(h.logger, w, twilioRes)

	default:
		twilioRes.Add(&twiml.Hangup{})
		writeTwilioResp(h.logger, w, twilioRes)
	}
}
