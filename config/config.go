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
}

// LoadConfig loads configuration from the environment
func LoadConfig() (*Config, error) {
	destNum := os.Getenv("DESTINATION_NUMBER")

	if len(destNum) == 0 {
		return nil, fmt.Errorf("DESTINATION_NUMBER environment variable must" +
			"be set")
	}

	return &Config{
		DestinationNumber: destNum,
	}, nil
}
