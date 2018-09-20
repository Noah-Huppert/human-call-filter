package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Noah-Huppert/human-call-filter/calls"
	"github.com/Noah-Huppert/human-call-filter/config"
	"github.com/Noah-Huppert/human-call-filter/libdb"

	"github.com/Noah-Huppert/golog"
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

	// Setup twilio call handler server
	callServer := calls.NewServer(logger, cfg, db)

	go func() {
		<-ctx.Done()

		err := callServer.Shutdown(context.Background())
		if err != nil {
			logger.Fatalf("failed to shutdown http call server: %s", err.Error())
		}
	}()

	logger.Debugf("starting http call server on %s", callServer.Addr)

	err = callServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error running http call server: %s", err.Error())
	}
}
