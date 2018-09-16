package calls

import (
	"net/http"

	"github.com/BTBurke/twiml"
	"github.com/Noah-Huppert/golog"
)

// writeTwilioResp writes a Twilio response as an HTTP response
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
