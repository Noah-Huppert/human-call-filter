package libdb

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

// Migrate runs migrations on a database
func Migrate(db *sql.DB) error {
	// Make database driver
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating database driver: %s", err.Error())
	}

	// Create migrate client
	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("error creating migrator: %s", err.Error())
	}

	// Migrate
	err = migrator.Up()
	if err != nil {
		return fmt.Errorf("error running migrations: %s", err.Error())
	}

	return nil
}
