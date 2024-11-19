package config

import (
	"log"
	"os"
	"zip/models"
)

var Config *models.Config

const (
	DefaultPort = 8000
	DefaultEnv  = "local"
)

// SetupConfig parses command-line arguments, environment variables from a YAML file, and configures the logger.
// If the help flag is provided, the program prints the usage information and exits with code 0.
// If an error occurs during flag validation, the program terminates with code 1.
func SetupConfig() *models.Config {
	conf := newConfig()

	if err := ParseFlags(conf); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	SetLogger(conf)

	conf.Email = os.Getenv("EMAIL")
	conf.Password = os.Getenv("PASSWORD")

	Config = conf
	return conf
}

func newConfig() *models.Config {
	return &models.Config{
		Port: DefaultPort,
		Env:  DefaultEnv,
	}
}
