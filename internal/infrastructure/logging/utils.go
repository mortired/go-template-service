package logging

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

// CheckElasticsearchHealthManually performs a manual health check of Elasticsearch
// This function can be called manually to verify Elasticsearch connectivity
func CheckElasticsearchHealthManually() error {
	// Load configuration from cache
	esConfig, err := GetElasticsearchConfig()
	if err != nil {
		return fmt.Errorf("failed to load Elasticsearch configuration: %w", err)
	}

	if esConfig == nil || !esConfig.Enabled {
		return fmt.Errorf("Elasticsearch is not configured")
	}

	// Perform health check
	return CheckElasticsearchHealth(*esConfig)
}

// GetElasticsearchStatus returns the current status of Elasticsearch configuration
func GetElasticsearchStatus() (map[string]interface{}, error) {
	esConfig, err := GetElasticsearchConfig()
	if err != nil {
		return nil, err
	}

	status := map[string]interface{}{
		"enabled":              esConfig != nil && esConfig.Enabled,
		"urls":                 nil,
		"index":                "",
		"health_check_enabled": false,
	}

	if esConfig != nil && esConfig.Enabled {
		status["urls"] = esConfig.URLs
		status["index"] = esConfig.Index
		status["health_check_enabled"] = esConfig.HealthCheckEnabled
		// Don't perform health check automatically, just show configuration
		status["health_status"] = "not_checked"
		status["note"] = "Use CheckElasticsearchHealthManually() to verify connectivity"
	}

	return status, nil
}

// TestElasticsearchConnection tests the connection to Elasticsearch without creating logs
func TestElasticsearchConnection() error {
	fmt.Println("🧪 Testing Elasticsearch connection...")

	esConfig, err := GetElasticsearchConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if esConfig == nil || !esConfig.Enabled {
		return fmt.Errorf("Elasticsearch is not configured")
	}

	// Perform a quick connection test
	return CheckElasticsearchHealth(*esConfig)
}

// GetElasticsearchClient creates and returns an Elasticsearch client for manual operations
func GetElasticsearchClient() (*elasticsearch.Client, error) {
	esConfig, err := GetElasticsearchConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if esConfig == nil || !esConfig.Enabled {
		return nil, fmt.Errorf("Elasticsearch is not configured")
	}

	return NewElasticsearchClient(*esConfig)
}
