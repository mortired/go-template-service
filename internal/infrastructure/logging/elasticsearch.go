package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"users/internal/infrastructure/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ElasticsearchConfig holds Elasticsearch configuration
type ElasticsearchConfig struct {
	URLs     []string
	Username string
	Password string
	Index    string
	Enabled  bool
	// Health check settings
	HealthCheckEnabled bool
	HealthCheckTimeout time.Duration
}

// Global configuration instance (singleton)
var (
	globalESConfig     *ElasticsearchConfig
	globalESConfigOnce sync.Once
	globalESConfigErr  error
)

// GetElasticsearchConfig returns the global Elasticsearch configuration
// This function loads configuration only once and caches it
func GetElasticsearchConfig() (*ElasticsearchConfig, error) {
	globalESConfigOnce.Do(func() {
		config, err := LoadElasticsearchConfigFromEnv()
		if err != nil {
			globalESConfigErr = err
			return
		}
		globalESConfig = &config
		globalESConfigErr = nil
	})
	return globalESConfig, globalESConfigErr
}

// ResetElasticsearchConfig resets the global configuration (mainly for testing)
func ResetElasticsearchConfig() {
	globalESConfigOnce = sync.Once{}
	globalESConfig = nil
	globalESConfigErr = nil
}

// ElasticsearchCore implements zapcore.Core for Elasticsearch logging
type ElasticsearchCore struct {
	esClient      *elasticsearch.Client
	index         string
	bulk          esutil.BulkIndexer
	healthChecked bool
	baseCore      zapcore.Core // Добавляем ссылку на базовый core для проверки уровня
}

// Enabled implements zapcore.Core
func (ec *ElasticsearchCore) Enabled(level zapcore.Level) bool {
	// Используем тот же уровень, что и базовый логгер
	if ec.baseCore != nil {
		return ec.baseCore.Enabled(level)
	}
	// Fallback: включен только для Info и выше
	return level >= zapcore.InfoLevel
}

// With implements zapcore.Core
func (ec *ElasticsearchCore) With(fields []zap.Field) zapcore.Core {
	return &ElasticsearchCore{
		esClient:      ec.esClient,
		index:         ec.index,
		bulk:          ec.bulk,
		healthChecked: ec.healthChecked,
		baseCore:      ec.baseCore,
	}
}

// Check implements zapcore.Core
func (ec *ElasticsearchCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if ec.Enabled(ent.Level) {
		return ce.AddCore(ent, ec)
	}
	return ce
}

// Sync implements zapcore.Core
func (ec *ElasticsearchCore) Sync() error {
	if ec.bulk != nil {
		// Add timeout to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := ec.bulk.Close(ctx); err != nil {
			fmt.Printf("⚠️  Failed to close bulk indexer: %v\n", err)
			return err
		}
	}
	return nil
}

// Write implements zapcore.Core
func (ec *ElasticsearchCore) Write(entry zapcore.Entry, fields []zap.Field) error {
	// If Elasticsearch is disabled, just return
	if ec.esClient == nil {
		return nil
	}

	// Send to Elasticsearch in background (non-blocking)
	go func() {
		// Add panic recovery to prevent crashes
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("⚠️  Panic in Elasticsearch logging goroutine: %v\n", r)
			}
		}()

		if err := ec.sendToElasticsearch(entry, fields); err != nil {
			// Log error to console but don't fail the main logging
			fmt.Printf("⚠️  Elasticsearch logging failed for %s: %v\n", entry.Message, err)
		}
	}()

	return nil
}

// sendToElasticsearch sends a log entry to Elasticsearch
func (ec *ElasticsearchCore) sendToElasticsearch(entry zapcore.Entry, fields []zap.Field) error {
	// Check if bulk indexer is available and working
	if ec.bulk != nil {
		// Try to use bulk indexer first
		doc := map[string]interface{}{
			"@timestamp": entry.Time,
			"level":      entry.Level.String(),
			"message":    entry.Message,
			"logger":     entry.LoggerName,
			"caller":     entry.Caller.String(),
			"stack":      entry.Stack,
		}

		// Add custom fields
		for _, field := range fields {
			switch field.Type {
			case zapcore.StringType:
				doc[field.Key] = field.String
			case zapcore.Int64Type:
				doc[field.Key] = field.Integer
			case zapcore.Float64Type:
				doc[field.Key] = field.Interface
			case zapcore.BoolType:
				doc[field.Key] = field.Integer == 1
			case zapcore.TimeType:
				doc[field.Key] = field.Interface
			default:
				doc[field.Key] = field.Interface
			}
		}

		// Try to add to bulk indexer, but handle errors gracefully
		jsonDoc, err := json.Marshal(doc)
		if err != nil {
			// Fall back to direct API call
			return ec.sendDirectToElasticsearch(entry, fields)
		}

		err = ec.bulk.Add(context.Background(), esutil.BulkIndexerItem{
			Index:  ec.index,
			Action: "index",
			Body:   bytes.NewReader(jsonDoc),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
				if config.GetEnvOrDefault("LOG_ELASTIC_DEBUG", "false") == "true" {
					fmt.Printf("⚠️  Bulk indexing failed for %s: %v\n", entry.Message, err)
				}
			},
		})

		// If bulk indexing fails, fall back to direct API call
		if err != nil {
			// Fall back to direct API call
			return ec.sendDirectToElasticsearch(entry, fields)
		}

		return nil
	}

	// If no bulk indexer, use direct API call
	return ec.sendDirectToElasticsearch(entry, fields)
}

// sendDirectToElasticsearch sends log entry to Elasticsearch directly via API
func (ec *ElasticsearchCore) sendDirectToElasticsearch(entry zapcore.Entry, fields []zap.Field) error {
	// Create document for Elasticsearch
	doc := map[string]interface{}{
		"@timestamp": entry.Time,
		"level":      entry.Level.String(),
		"message":    entry.Message,
		"logger":     entry.LoggerName,
		"caller":     entry.Caller.String(),
		"stack":      entry.Stack,
	}

	// Add custom fields
	for _, field := range fields {
		switch field.Type {
		case zapcore.StringType:
			doc[field.Key] = field.String
		case zapcore.Int64Type:
			doc[field.Key] = field.Integer
		case zapcore.Float64Type:
			doc[field.Key] = field.Interface
		case zapcore.BoolType:
			doc[field.Key] = field.Integer == 1
		case zapcore.TimeType:
			doc[field.Key] = field.Interface
		default:
			doc[field.Key] = field.Interface
		}
	}

	// Send document to Elasticsearch
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	resp, err := ec.esClient.Index(
		ec.index,
		bytes.NewReader(jsonData),
		ec.esClient.Index.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("failed to send document to Elasticsearch: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("Elasticsearch returned error: %s", resp.Status())
	}

	return nil
}

// Close closes the bulk indexer
func (ec *ElasticsearchCore) Close() error {
	if ec.bulk != nil {
		return ec.bulk.Close(context.Background())
	}
	return nil
}

// LoadElasticsearchConfigFromEnv loads Elasticsearch configuration from environment variables
func LoadElasticsearchConfigFromEnv() (ElasticsearchConfig, error) {
	url := config.GetEnvOrDefault("LOG_ELASTIC_URL", "")
	username := config.GetEnvOrDefault("LOG_ELASTIC_USERNAME", "")
	password := config.GetEnvOrDefault("LOG_ELASTIC_PASSWORD", "")
	index := config.GetEnvOrDefault("LOG_ELASTIC_INDEX", "")
	healthCheckEnabled := config.GetEnvAsBoolOrDefault("LOG_ELASTIC_HEALTH_CHECK", true)

	enabled := url != ""

	if !enabled && (username != "" || password != "" || index != "") {
		return ElasticsearchConfig{}, fmt.Errorf("LOG_ELASTIC_URL is required for Elasticsearch integration. Other variables will be ignored")
	}

	// If no Elasticsearch configuration at all, return disabled config without error
	if !enabled {
		return ElasticsearchConfig{
			URLs:               nil,
			Username:           "",
			Password:           "",
			Index:              "",
			Enabled:            false,
			HealthCheckEnabled: false,
			HealthCheckTimeout: 5 * time.Second,
		}, nil
	}

	config := ElasticsearchConfig{
		URLs:               []string{url},
		Username:           username,
		Password:           password,
		Index:              index,
		Enabled:            enabled,
		HealthCheckEnabled: healthCheckEnabled,
		HealthCheckTimeout: 5 * time.Second,
	}

	// Configuration info is now shown only once when GetElasticsearchConfig is called
	return config, nil
}

// NewElasticsearchClient creates a new Elasticsearch client
func NewElasticsearchClient(config ElasticsearchConfig) (*elasticsearch.Client, error) {
	if !config.Enabled {
		return nil, fmt.Errorf("Elasticsearch is disabled")
	}

	// Create Elasticsearch client configuration
	cfg := elasticsearch.Config{
		Addresses: config.URLs,
		Username:  config.Username,
		Password:  config.Password,
		// Add timeout settings
		Transport: &http.Transport{
			ResponseHeaderTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// Create client
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	return client, nil
}

// NewElasticsearchCore creates a new Elasticsearch core
func NewElasticsearchCore(baseCore zapcore.Core, esClient *elasticsearch.Client, index string) *ElasticsearchCore {
	// Create bulk indexer for better performance
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         index,
		Client:        esClient,
		NumWorkers:    1,                // Reduce workers to prevent conflicts
		FlushInterval: 10 * time.Second, // Flush more frequently
		OnError: func(ctx context.Context, err error) {
			fmt.Printf("⚠️  Bulk indexing error: %v\n", err)
		},
		OnFlushStart: func(ctx context.Context) context.Context {
			if config.GetEnvOrDefault("LOG_ELASTIC_DEBUG", "false") == "true" {
				fmt.Printf("📤 Starting bulk flush to Elasticsearch\n")
			}
			return ctx
		},
		OnFlushEnd: func(ctx context.Context) {
			if config.GetEnvOrDefault("LOG_ELASTIC_DEBUG", "false") == "true" {
				fmt.Printf("✅ Bulk flush completed\n")
			}
		},
	})

	if err != nil {
		fmt.Printf("⚠️  Failed to create bulk indexer: %v\n", err)
		// Continue without bulk indexing
		bulkIndexer = nil
	}

	esCore := &ElasticsearchCore{
		esClient:      esClient,
		index:         index,
		bulk:          bulkIndexer,
		healthChecked: false,
		baseCore:      baseCore,
	}

	return esCore
}

// NewLoggerWithElasticsearch creates a new logger with Elasticsearch integration
func NewLoggerWithElasticsearch(config Config, esConfig ElasticsearchConfig) (*Logger, error) {
	// Create base logger directly (avoid recursion)
	logger, err := NewStandardLogger(config)
	if err != nil {
		return nil, err
	}

	// If Elasticsearch is not enabled, return base logger
	if !esConfig.Enabled {
		return logger, nil
	}

	// Create Elasticsearch client
	esClient, err := NewElasticsearchClient(esConfig)
	if err != nil {
		fmt.Printf("⚠️  Failed to create Elasticsearch client: %v\n", err)
		// Return base logger without Elasticsearch
		return logger, nil
	}

	// Create Elasticsearch core
	esCore := NewElasticsearchCore(logger.Core(), esClient, esConfig.Index)

	// Create new logger with BOTH cores using zapcore.NewCore
	// This ensures that logs go to BOTH console AND Elasticsearch
	combinedCore := zapcore.NewTee(logger.Core(), esCore)

	esLogger := zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{esLogger}, nil
}

// CheckElasticsearchHealth checks if Elasticsearch is accessible and can write to the specified index
// This function is now separate and should be called manually when needed
func CheckElasticsearchHealth(config ElasticsearchConfig) error {
	if !config.Enabled {
		return fmt.Errorf("Elasticsearch is disabled")
	}

	fmt.Printf("🔍 Checking Elasticsearch health at %v...\n", config.URLs)

	// Create client for health check
	esClient, err := NewElasticsearchClient(config)
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// 1. Check cluster health
	resp, err := esClient.Cluster.Health(esClient.Cluster.Health.WithContext(context.Background()))
	if err != nil {
		return fmt.Errorf("failed to check cluster health: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("Elasticsearch cluster health check failed with status %d", resp.StatusCode)
	}

	fmt.Printf("✅ Cluster health check passed\n")

	// 2. Check if we can access the index (create if doesn't exist)
	indexExists, err := esClient.Indices.Exists([]string{config.Index})
	if err != nil {
		return fmt.Errorf("failed to check index %s: %w", config.Index, err)
	}

	// If index doesn't exist (404), try to create it
	if indexExists.StatusCode == 404 {
		fmt.Printf("📝 Index %s doesn't exist, creating...\n", config.Index)
		if err := createElasticsearchIndex(esClient, config); err != nil {
			return fmt.Errorf("failed to create index %s: %w", config.Index, err)
		}
		fmt.Printf("✅ Index %s created successfully\n", config.Index)
	} else if indexExists.StatusCode >= 400 {
		return fmt.Errorf("failed to access index %s with status %d", config.Index, indexExists.StatusCode)
	} else {
		fmt.Printf("✅ Index %s exists and accessible\n", config.Index)
	}

	// 3. Test write permissions with a test document
	fmt.Printf("🧪 Testing write permissions...\n")
	if err := testElasticsearchWrite(esClient, config); err != nil {
		return fmt.Errorf("failed to test write permissions to index %s: %w", config.Index, err)
	}

	fmt.Printf("✅ Elasticsearch at %v is accessible and ready for logging\n", config.URLs)
	return nil
}

// createElasticsearchIndex creates the Elasticsearch index if it doesn't exist
func createElasticsearchIndex(client *elasticsearch.Client, config ElasticsearchConfig) error {
	indexMapping := `{
		"mappings": {
			"properties": {
				"@timestamp": {
					"type": "date"
				},
				"level": {
					"type": "keyword"
				},
				"message": {
					"type": "text"
				},
				"logger": {
					"type": "keyword"
				},
				"service": {
					"type": "keyword"
				},
				"version": {
					"type": "keyword"
				},
				"environment": {
					"type": "keyword"
				}
			}
		}
	}`

	resp, err := client.Indices.Create(
		config.Index,
		client.Indices.Create.WithBody(strings.NewReader(indexMapping)),
		client.Indices.Create.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("failed to create index with status %d", resp.StatusCode)
	}

	return nil
}

// testElasticsearchWrite tests if we can write to the Elasticsearch index
func testElasticsearchWrite(client *elasticsearch.Client, config ElasticsearchConfig) error {
	testDoc := map[string]interface{}{
		"@timestamp":  time.Now().Format(time.RFC3339),
		"level":       "INFO",
		"message":     "Elasticsearch connection test",
		"logger":      "health-check",
		"service":     "users",
		"environment": "test",
	}

	jsonData, err := json.Marshal(testDoc)
	if err != nil {
		return err
	}

	resp, err := client.Index(
		config.Index,
		bytes.NewReader(jsonData),
		client.Index.WithDocumentID("health-check-test"),
		client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("failed to write test document with status %d", resp.StatusCode)
	}

	// Clean up test document
	client.Delete(config.Index, "health-check-test", client.Delete.WithContext(context.Background()))

	return nil
}
