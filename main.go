package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Noah-Huppert/human-call-filter/calls"
	"github.com/Noah-Huppert/human-call-filter/config"

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

	// Load application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("error loading configuration: %s", err.Error())
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Setup twilio number handler
	router := mux.NewRouter()
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	router.Handle("/call", calls.NewCallsHandler(logger)).Methods("POST")
	router.Handle("/input/test/{eq}", calls.NewTestInputHandler(logger, cfg)).
		Methods("POST")
	router.Handle("/audio-clips/{file}.mp3", calls.NewAudioClipsHandler(logger))

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
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

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error running http server: %s", err.Error())
	}
}
