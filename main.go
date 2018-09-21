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
	"github.com/Noah-Huppert/human-call-filter/dashboard"
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

	// waitChan is used to wait for http servers to shut down, a nil on the
	// channel is a success. Otherwise the shutdown error is sent.
	waitChan := make(chan WaitEntry)

	totalWaitEntries := 2
	receivedWaitEntries := 0

	// Setup twilio call handler server
	callServer := calls.NewServer(logger, cfg, db)

	shutdownHTTPServerOnExit("calls", ctx, waitChan, &callServer)

	logger.Debugf("starting calls http server on %s", callServer.Addr)
	startHTTPServer("calls", waitChan, &callServer)

	// Setup dashboard server
	dashboardServer := dashboard.NewServer(logger, cfg, db)

	shutdownHTTPServerOnExit("dashboard", ctx, waitChan, &dashboardServer)

	logger.Debugf("starting dashboard http server on %s", dashboardServer.Addr)
	startHTTPServer("dashboard", waitChan, &dashboardServer)

	// Wait for servers to exit
	ctxDone := false

	for receivedWaitEntries < totalWaitEntries {
		select {
		case <-ctx.Done():
			ctxDone = true

		case waitEntry := <-waitChan:
			if waitEntry.err != nil {
				logger.Errorf("error running %s: %s", waitEntry.name, err.Error())

				// Cancel context if not canceled already, so that other jobs
				// shut down
				if !ctxDone {
					cancelFn()
				}
			} else {
				logger.Debugf("%s successfully completed", waitEntry.name)
			}

			receivedWaitEntries++
		}
	}
}

// WaitEntry holds information about the status of an asynchronous job
type WaitEntry struct {
	// name identifies the job
	name string

	// err holds an error if one occurs while running the job, nil if the job
	// completed successfully
	err error
}

// startHTTPServer starts an http server in a go routine and exits if there is
// an error.
//
// The name argument is used to identify the http server if an error occurs.
func startHTTPServer(name string, waitChan chan<- WaitEntry, server *http.Server) {
	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			waitChan <- WaitEntry{
				name: name,
				err:  fmt.Errorf("error running http server: %s", err.Error()),
			}
		}
	}()
}

// shutdownHTTPServerOnExit starts a go routine which waits for a Context to
// close and then gracefully shuts down an http.Server.
//
// The name argument will be used to identify the http server if an error
// occurs.
func shutdownHTTPServerOnExit(name string, ctx context.Context,
	waitChan chan<- WaitEntry, server *http.Server) {

	go func() {
		<-ctx.Done()

		err := server.Shutdown(context.Background())
		if err != nil {
			waitChan <- WaitEntry{
				name: name,
				err: fmt.Errorf("failed to shutdown the http server: %s",
					err.Error()),
			}
		}

		waitChan <- WaitEntry{
			name: name,
			err:  nil,
		}
	}()
}
