package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Noah-Huppert/golog"

	"github.com/BTBurke/twiml"
)

func main() {
	// Context
	ctx, cancelFn := context.WithCancel(context.Background())

	// Logger
	logger := golog.NewStdLogger("spam-call-blocker")

	// Setup exit handler
	exitSigChan := make(chan os.Signal, 1)
	signal.Notify(exitSigChan, os.Interrupt)

	go func() {
		<-exitSigChan
		cancelFn()
	}()

	// Setup twilio number handler
	server := http.Server{
		Addr: ":8000",
		Handler: CallHandler{
			logger: logger,
		},
	}

	go func() {
		<-ctx.Done()

		err := server.Shutdown(context.Background())
		if err != nil {
			logger.Fatalf("failed to shutdown http server: %s", err.Error())
		}
	}()

	logger.Debugf("starting http server on %s", server.Addr)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error running http server: %s", err.Error())
	}
}

// CallHandler responds to Twilio phone calls
type CallHandler struct {
	logger golog.Logger
}

// writeXML writes a Twilio response as an HTTP response
func writeTwilioResp(logger golog.Logger, w http.ResponseWriter, twilioRes *twiml.Response) {
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
	_, err = w.Write(bytes)
	if err != nil {
		logger.Errorf("error writing twilio response: %s", err.Error())

		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
}

// ServeHTTP handles a Twilio phone call
func (h CallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		return

	case twiml.Ringing, twiml.Queued:
		twilioRes.Add(&twiml.Say{
			Text: "hello world",
		})
		writeTwilioResp(h.logger, w, twilioRes)

	default:
		twilioRes.Add(&twiml.Hangup{})
		writeTwilioResp(h.logger, w, twilioRes)
	}
}
