package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Noah-Huppert/human-call-filter/config"
	"github.com/Noah-Huppert/human-call-filter/handlers"
	"github.com/Noah-Huppert/human-call-filter/libdb"

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

	// Connect to database
	db, err := libdb.ConnectX(cfg.DBConfig)
	if err != nil {
		logger.Fatalf("error connecting to database: %s", err.Error())
	}

	// Setup twilio number handler
	router := mux.NewRouter()

	routesLogger := logger.GetChild("routes")

	router.Handle("/healthz", handlers.NewHealthHandler())

	router.Handle("/call", handlers.NewCallsHandler(
		routesLogger.GetChild("call"), cfg, db)).Methods("POST")

	router.Handle("/input/challenge/{challenge_id}",
		handlers.NewTestInputHandler(routesLogger.GetChild("input"), cfg,
			db)).Methods("POST")

	router.Handle("/audio-clips/{file}.mp3", handlers.NewAudioClipsHandler(logger))

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
