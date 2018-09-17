package handlers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Noah-Huppert/human-call-filter/models"

	"github.com/BTBurke/twiml"
	"github.com/Noah-Huppert/golog"
	"github.com/jmoiron/sqlx"
)

// CallsHandler implements a http.Handler which responds to Twilio phone calls
type CallsHandler struct {
	// logger records debug information
	logger golog.Logger

	// db is a database connection
	db *sqlx.DB
}

// NewCallsHandler creates a new CallsHandler
func NewCallsHandler(logger golog.Logger, db *sqlx.DB) CallsHandler {
	return CallsHandler{
		logger: logger,
		db:     db,
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

	// Query database for phone number
	phoneNum := &models.PhoneNumber{
		Number: twilioReq.From,
	}

	err = phoneNum.QueryByNumber(h.db)
	if err == sql.ErrNoRows {
		h.logger.Debugf("No phone number found for: %#v", phoneNum)
	} else if err != nil {
		h.logger.Errorf("error querying database for phone number: %s",
			err.Error())
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
