package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds application configuration
type Config struct {
	// DBConfig holds database configuration
	DBConfig

	// DestinationNumber is the phone number callers will be forwarded to if
	// they pass the challenge
	DestinationNumber string `required:"true" envconfig:"destination_number"`

	// CallsHTTPPort is the port the server which responds to twilio calls
	// will run on
	CallsHTTPPort string `default:"8000" envconfig:"calls_http_port"`

	// DashboardHTTPPort is the port the server which serves an internal
	// dashboard will run on
	DashboardHTTPPort string `default:"80001" envconfig:"dashboard_http_port"`
}

// LoadConfig loads configuration from the environment
func LoadConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error loading configuration from environment"+
			" variables: %s", err.Error())
	}

	return &cfg, nil
}
