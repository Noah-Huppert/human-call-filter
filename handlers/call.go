package handlers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

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
	// Note time request was received for later use in PhoneCall model
	callReceived := time.Now()

	// Parse request
	var twilioReq twiml.VoiceRequest

	err := twiml.Bind(&twilioReq, r)
	if err != nil {
		h.logger.Errorf("error reading Twilio request: %s", err.Error())

		writeStatus(w, http.StatusBadRequest)
		return
	}

	// Query database for phone number
	phoneNum := &models.PhoneNumber{
		Number:  twilioReq.From,
		Name:    twilioReq.CallerName,
		State:   twilioReq.FromState,
		City:    twilioReq.FromCity,
		ZipCode: twilioReq.FromZip,
	}

	err = phoneNum.QueryByNumber(h.db)

	// No number found, insert
	if err == sql.ErrNoRows {
		err = phoneNum.Insert(h.db)
		if err != nil {
			h.logger.Errorf("error inserting new phone number into db: %s",
				err.Error())

			writeStatus(w, http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		h.logger.Errorf("error querying database for phone number: %s",
			err.Error())

		writeStatus(w, http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("number: %#v", phoneNum)

	// Insert phone call row into database
	phoneCall := &models.PhoneCall{
		PhoneNumberID: phoneNum.ID,
		DateReceived:  callReceived,
	}

	err = phoneCall.Insert(h.db)
	if err != nil {
		h.logger.Errorf("error inserting new phone call into db: %s",
			err.Error())

		writeStatus(w, http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("phone call: %#v", phoneCall)

	// Create response
	twilioRes := twiml.NewResponse()

	h.logger.Debugf("received call request %#v", twilioReq)

	switch twilioReq.CallStatus {
	case twiml.InProgress:
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, http.StatusText(http.StatusOK))

	case twiml.Ringing, twiml.Queued:
		// Generate challenge
		a := rand.Intn(3) + 1
		b := rand.Intn(3) + 1
		eq := a + b

		// Insert challenge into db
		challenge := &models.Challenge{
			PhoneCallID: phoneCall.ID,
			DateAsked:   time.Now(),
			OperandA:    a,
			OperandB:    b,
			Solution:    eq,
			Status:      models.ChallengeStatusAnswering,
		}

		err = challenge.Insert(h.db)
		if err != nil {
			h.logger.Errorf("error inserting challenge into db: %s",
				err.Error())

			writeStatus(w, http.StatusInternalServerError)
			return
		}

		h.logger.Debugf("challenge: %#v", challenge)

		// Ask question to user
		twilioRes.Add(&twiml.Play{
			URL:    "/audio-clips/intro.mp3",
			Digits: "w",
		})
		twilioRes.Add(&twiml.Say{
			Voice: "man",
			Text:  fmt.Sprintf("What is ,%d. plus ,%d", a, b),
		})
		twilioRes.Add(&twiml.Gather{
			Action:    fmt.Sprintf("/input/challenge/%d", challenge.ID),
			NumDigits: 1,
			Timeout:   10,
		})
		writeTwilioResp(h.logger, w, twilioRes)

	default:
		twilioRes.Add(&twiml.Hangup{})
		writeTwilioResp(h.logger, w, twilioRes)
	}
}
