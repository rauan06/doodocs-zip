package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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

	if err := loadEnvFromYAML("internal/config/config.yaml", conf); err != nil {
		log.Fatalf("Error loading environment variables from YAML: %v", err)
	}

	if err := ParseFlags(conf); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	SetLogger(conf)

	Config = conf
	return conf
}

func newConfig() *models.Config {
	return &models.Config{
		Port: DefaultPort,
		Env:  DefaultEnv,
	}
}

// loadEnvFromYAML loads environment variables from a YAML file and updates the Config struct
func loadEnvFromYAML(filePath string, conf *models.Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var isEnvSection bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check if we are in the "env_variables" section
		if strings.HasPrefix(line, "env_variables:") {
			isEnvSection = true
			continue
		}

		if isEnvSection && line != "" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "EMAIL":
				conf.Email = value
			case "PASSWORD":
				conf.Password = value
			}
		}
	}

	// If there was an error scanning the file
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}
