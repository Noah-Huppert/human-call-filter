package main

import (
	"github.com/Noah-Huppert/human-call-filter/config"
	"github.com/Noah-Huppert/human-call-filter/libdb"

	"github.com/Noah-Huppert/golog"
)

func main() {
	// Setup logger
	logger := golog.NewStdLogger("db-migrate")

	// Load configuration
	cfg, err := config.NewDBConfig()
	if err != nil {
		logger.Fatalf("failed to load database configuration: %s", err.Error())
	}

	// Connect to database
	db, err := libdb.Connect(*cfg)
	if err != nil {
		logger.Fatalf("failed to connect to database: %s", err.Error())
	}

	// Run migrations
	logger.Info("running migrations")

	err = libdb.Migrate(db)
	if err != nil {
		logger.Fatalf("error running migrations: %s", err.Error())
	}

	logger.Info("successfully ran migrations")
}
