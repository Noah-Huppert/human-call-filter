package handlers

import (
	"net/http"

	"github.com/BTBurke/twiml"
	"github.com/Noah-Huppert/golog"
)

// writeStatus responds with a status code and the related status text
func writeStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// writeTwilioResp writes a Twilio response as an HTTP response
func writeTwilioResp(logger golog.Logger, w http.ResponseWriter,
	twilioRes *twiml.Response) {

	// Encode Twilio response to bytes
	bytes, err := twilioRes.Encode()

	if err != nil {
		logger.Errorf("error encoding twilio response into bytes: %s",
			err.Error())

		writeStatus(http.StatusInternalServerError)
		return
	}

	// Write bytes as response
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bytes)
	if err != nil {
		logger.Errorf("error writing twilio response: %s", err.Error())

		writeStatus(http.StatusInternalServerError)
		return
	}
}
