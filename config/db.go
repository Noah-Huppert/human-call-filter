package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// DBConfig holds database related application configuration
type DBConfig struct {
	// DBHost is the host of the database
	// Set via the APP_DB_HOST environment variable
	DBHost string `required:"true" envconfig:"db_host"`

	// DBPort is the port the database accepts connections on, defaults to 5432
	// Set via the APP_DB_PORT environment variable
	DBPort int `default:"5432" envconfig:"db_port"`

	// DBName is the name of the database to save data in
	// Set via the APP_DB_NAME environment variable
	DBName string `required:"true" envconfig:"db_name"`

	// DBUsername is the username used to authenticate with the database
	// set via the APP_DB_USERNAME environment variable
	DBUsername string `required:"true" envconfig:"db_username"`

	// DBPassword is the password used to authenticate with the database
	// Set via the APP_DB_PASSWORD environment variable
	DBPassword string `envconfig:"db_password"`
}

// NewDBConfig loads database Config values from environment variables. Variables names
// will be capitalized field names from the Config struct, prefixed
// with APP.
func NewDBConfig() (*DBConfig, error) {
	var cfg DBConfig

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading database configuration from the"+
			" environment: %s", err.Error())
	}

	return &cfg, nil
}
