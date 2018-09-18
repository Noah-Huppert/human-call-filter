package libdb

import (
	"database/sql"
	"fmt"

	"github.com/Noah-Huppert/human-call-filter/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect to a PostgreSQL database
func Connect(dbCfg config.DBConfig) (*sql.DB, error) {
	// Assemble database connection string
	sqlConnStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s "+
		"sslmode=disable", dbCfg.DBHost, dbCfg.DBPort,
		dbCfg.DBName, dbCfg.DBUsername)

	if len(dbCfg.DBPassword) > 0 {
		sqlConnStr += fmt.Sprintf(" password=%s", dbCfg.DBPassword)
	}

	// Connect to database
	db, err := sql.Open("postgres", sqlConnStr)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to database: %s",
			err.Error())
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error checking connection: %s", err.Error())
	}

	return db, nil
}

// ConnectX connects to a PostgreSQL database an then makes an sqlx.DB instance
func ConnectX(dbCfg config.DBConfig) (*sqlx.DB, error) {
	// Connect
	db, err := Connect(dbCfg)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err.Error())
	}

	// Make sqlx.DB instance
	dbx := sqlx.NewDb(db, "postgres")

	return dbx, nil
}
