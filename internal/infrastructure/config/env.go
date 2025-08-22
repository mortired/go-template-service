package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	envLoaded bool
	envMutex  sync.Once
)

// loadEnvOnce loads environment variables from .env file if it exists
// This function is called only once thanks to sync.Once
func loadEnvOnce() error {
	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// File doesn't exist, use only system environment variables
		return nil
	}

	// Open .env file
	file, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE line
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"'`)

		// Set environment variable only if it's not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}

// LoadEnv loads environment variables from .env file if it exists
// This function is kept for backward compatibility
func LoadEnv() error {
	return ensureEnvLoaded()
}

// ensureEnvLoaded ensures that .env file is loaded only once
func ensureEnvLoaded() error {
	var loadErr error
	envMutex.Do(func() {
		loadErr = loadEnvOnce()
		if loadErr != nil {
			// Use simple fmt.Printf to avoid import cycle
			fmt.Printf("⚠️  Warning: failed to load .env file: %v\n", loadErr)
		}
		envLoaded = true
	})
	return loadErr
}

// GetEnvOrDefault returns environment variable value or default value
func GetEnvOrDefault(key, defaultValue string) string {
	// Ensure .env file is loaded before reading environment variables
	ensureEnvLoaded()

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsIntOrDefault returns environment variable value as int or default value
func GetEnvAsIntOrDefault(key string, defaultValue int) int {
	// Ensure .env file is loaded before reading environment variables
	ensureEnvLoaded()

	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsBoolOrDefault returns environment variable value as bool or default value
func GetEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	// Ensure .env file is loaded before reading environment variables
	ensureEnvLoaded()

	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(strings.ToLower(value)); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetEnvRequired returns environment variable value or error if not set
func GetEnvRequired(key string) (string, error) {
	// Ensure .env file is loaded before reading environment variables
	ensureEnvLoaded()

	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("required environment variable %s is not set", key)
}

// GetEnvAsIntRequired returns environment variable value as int or error if not set
func GetEnvAsIntRequired(key string) (int, error) {
	// Ensure .env file is loaded before reading environment variables
	ensureEnvLoaded()

	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue, nil
		}
		return 0, fmt.Errorf("environment variable %s is not a valid integer", key)
	}
	return 0, fmt.Errorf("required environment variable %s is not set", key)
}
