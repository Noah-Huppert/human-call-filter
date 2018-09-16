package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	// DestinationNumber is the phone number callers will be forwarded to if
	// they pass the challenge
	DestinationNumber string

	// HTTPPort is the port to serve http traffic on, defaults to 8000
	HTTPPort string
}

// LoadConfig loads configuration from the environment
func LoadConfig() (*Config, error) {
	// Destination number
	destNum := os.Getenv("DESTINATION_NUMBER")

	if len(destNum) == 0 {
		return nil, fmt.Errorf("DESTINATION_NUMBER environment variable must" +
			"be set")
	}

	// HTTP port
	httpPort := os.Getenv("HTTP_PORT")

	if len(httpPort) == 0 {
		httpPort = "8000"
	}

	return &Config{
		DestinationNumber: destNum,
		HTTPPort:          httpPort,
	}, nil
}
