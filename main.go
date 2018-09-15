package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/Noah-Huppert/human-call-filter/calls"

	"github.com/Noah-Huppert/golog"
	"github.com/gorilla/mux"
)

func main() {
	// Context
	ctx, cancelFn := context.WithCancel(context.Background())

	// Logger
	logger := golog.NewStdLogger("human-call-filter")

	// Setup exit handler
	exitSigChan := make(chan os.Signal, 1)
	signal.Notify(exitSigChan, os.Interrupt)

	go func() {
		<-exitSigChan
		cancelFn()
	}()

	// Setup twilio number handler
	router := mux.NewRouter()
	router.Handle("/call", calls.NewCallsHandler(logger)).Methods("POST")
	router.PathPrefix("/audio-clips").Handler(http.StripPrefix("/audio-clips/",
		http.FileServer(http.Dir("./audio-clips"))))

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
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
