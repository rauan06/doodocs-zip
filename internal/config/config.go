package config

import (
	"log"
	"zip/models"
)

const (
	DefaultPort = 8000
	DefaultEnv  = "local"
)

// SetupConfig parses command-line arguments using the flag package and configures the logger.
// If the help flag is provided, the program prints the usage information and exits with code 0.
// If an error occurs during flag validation, the program terminates with code 1.
func SetupConfig() *models.Config {
	conf := newConfig()

	if err := ParseFlags(conf); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	SetLogger(conf)

	return conf
}

func newConfig() *models.Config {
	return &models.Config{
		Port: DefaultPort,
		Env:  DefaultEnv,
	}
}
