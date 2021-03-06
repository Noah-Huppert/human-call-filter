package calls

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Noah-Huppert/human-call-filter/config"
	"github.com/Noah-Huppert/human-call-filter/models"

	"github.com/BTBurke/twiml"
	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// TestInputHandler is an http.Handler which receives digits provided by
// callers when they attempt to answer the challenge question
type TestInputHandler struct {
	// logger is used to record debug information
	logger golog.Logger

	// cfg holds application configuration
	cfg *config.Config

	// db is a database connection
	db *sqlx.DB
}

// NewTestInputHandler creates a TestInputHandler
func NewTestInputHandler(logger golog.Logger, cfg *config.Config,
	db *sqlx.DB) TestInputHandler {

	return TestInputHandler{
		logger: logger,
		cfg:    cfg,
		db:     db,
	}
}

// ServeHTTP implements http.Handler
func (h TestInputHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get challenge id from URL
	vars := mux.Vars(r)

	challengeID, err := strconv.ParseInt(vars["challenge_id"], 10, 64)
	if err != nil {
		h.logger.Errorf("error parsing challenge id URL parameter into "+
			"int: %s", err.Error())

		writeStatus(w, http.StatusBadRequest)
		return
	}

	// Query database for challenge id
	challenge := &models.Challenge{
		ID: challengeID,
	}

	err = challenge.QueryByIDForAnswer(h.db)

	if err == sql.ErrNoRows {
		h.logger.Errorf("received input for challenge which didn't exist, "+
			"challenge id: %d", challenge.ID)

		writeStatus(w, http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Errorf("error query database for challenge: %s", err.Error())

		writeStatus(w, http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("challenge: %#v", challenge)

	// Check challenge hasn't already been answered
	if challenge.Status != models.ChallengeStatusAnswering {
		h.logger.Error("received input for challenge which has already " +
			"been answered")

		writeStatus(w, http.StatusBadRequest)
		return
	}

	// Parse request
	var twilioReq twiml.RecordActionRequest

	err = twiml.Bind(&twilioReq, r)
	if err != nil {
		h.logger.Errorf("error reading Twilio request: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	h.logger.Debugf("received input request: %s", twilioReq)

	// Get challenge response
	enteredStr := twilioReq.Digits

	entered, err := strconv.Atoi(enteredStr)
	if err != nil {
		h.logger.Errorf("error parsing entered digits into int, digits: %s,"+
			" err: %s", enteredStr, err.Error())

		writeStatus(w, http.StatusInternalServerError)
		return
	}

	// Check challenge solution
	if challenge.Solution == entered {
		challenge.Status = models.ChallengeStatusPassed
	} else {
		challenge.Status = models.ChallengeStatusFailed
	}

	// Make twilio response based on challenge success or failure
	twilioRes := twiml.NewResponse()

	if challenge.Status == models.ChallengeStatusPassed {
		twilioRes.Add(&twiml.Play{
			URL:    "/audio-clips/success.mp3",
			Digits: "w",
		})

		dial := &twiml.Dial{}

		dial.Add(&twiml.Sip{
			Address: h.cfg.Destination,
		})

		twilioRes.Add(dial)

		h.logger.Debug("caller passed challenge")
	} else {
		twilioRes.Add(&twiml.Play{
			URL:    "/audio-clips/fail.mp3",
			Digits: "w",
		})

		h.logger.Debug("caller failed challenge")
	}

	writeTwilioResp(h.logger, w, twilioRes)

	// Save challenge status
	err = challenge.UpdateStatusByID(h.db)
	if err != nil {
		h.logger.Errorf("error updating challenge status: %s", err.Error())
	}
}
