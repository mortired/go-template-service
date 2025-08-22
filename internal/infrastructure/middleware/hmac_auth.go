package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strings"
	"time"
	"users/internal/infrastructure/logging"
	"users/internal/infrastructure/response"

	"go.uber.org/zap"
)

// HMACClientSecret represents client secret
type HMACClientSecret struct {
	ClientID   string `json:"clientid"`
	Secret     string `json:"secret"`
	Department string `json:"department"`
	Descr      string `json:"descr"`
}

// HMACRouteRights represents route access rights for clients
type HMACRouteRights map[string][]string

// HMACConfig interface for HMAC configuration
type HMACConfig interface {
	GetClientSecrets() []HMACClientSecret
	GetRouteRights() HMACRouteRights
	GetAlgorithm() string
	GetMaxAge() int
	IsRequired() bool
}

// HMACAuthMiddleware creates HMAC authentication middleware
func HMACAuthMiddleware(cfg HMACConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create logger for this function
			logger, err := logging.New(logging.LoadConfigFromEnv())
			if err != nil {
				logger = nil
			}

			// Validate HMAC signature
			if err := validateHMACSignature(r, cfg); err != nil {
				if cfg.IsRequired() {
					// Create problem response
					problem := response.NewProblem(http.StatusUnauthorized, "Invalid HMAC signature").
						WithType(response.TypeUnauthorized).
						WithInstance(r.URL.Path).
						WithDetail(err.Error())

					// Send JSON response
					w.Header().Set("Content-Type", "application/problem+json")
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(problem)
					return
				}
				// If verification is not required, log error but continue
				if logger != nil {
					logger.Warn("HMAC validation warning", zap.String("error", err.Error()))
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// validateHMACSignature verifies HMAC signature of request
func validateHMACSignature(r *http.Request, cfg HMACConfig) error {
	// Get Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("missing Authorization header")
	}

	// Parse Authorization header
	credential, signature, err := parseAuthorizationHeader(authHeader)
	if err != nil {
		return fmt.Errorf("invalid Authorization header: %v", err)
	}

	// Get Date header
	dateHeader := r.Header.Get("Date")
	if dateHeader == "" {
		return fmt.Errorf("missing Date header")
	}

	// Get x-content-hmac header
	contentHash := r.Header.Get("x-content-hmac")
	if contentHash == "" {
		return fmt.Errorf("missing x-content-hmac header")
	}

	// Find client secret
	clientSecret, err := findClientSecret(credential, cfg.GetClientSecrets())
	if err != nil {
		return fmt.Errorf("client not found or invalid: %v", err)
	}

	// Check route access rights
	if err := checkRouteAccess(credential, r.URL.Path, cfg.GetRouteRights()); err != nil {
		return fmt.Errorf("access denied to route: %v", err)
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %v", err)
	}

	// Restore request body for further use
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Create string to sign in JavaScript script format
	stringToSign := createStringToSign(r, body, dateHeader, contentHash)

	// Verify signature
	algorithm := getAlgorithm(cfg.GetAlgorithm())
	if ok, err := verifySignature(stringToSign, signature, clientSecret, algorithm); !ok {
		if err != nil {
			return fmt.Errorf("invalid HMAC signature: %v", err)
		}
		return fmt.Errorf("invalid HMAC signature")
	}

	return nil
}

// parseAuthorizationHeader parses Authorization header
func parseAuthorizationHeader(authHeader string) (credential, signature string, err error) {
	// Format: HMAC-SHA256 Credential=example&SignedHeaders=date;host;x-content-hmac&Signature=...
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid Authorization header format")
	}

	// Extract parameters
	params := strings.Split(parts[1], "&")
	for _, param := range params {
		if strings.HasPrefix(param, "Credential=") {
			credential = strings.TrimPrefix(param, "Credential=")
		} else if strings.HasPrefix(param, "Signature=") {
			signature = strings.TrimPrefix(param, "Signature=")
		}
	}

	if credential == "" || signature == "" {
		return "", "", fmt.Errorf("missing Credential or Signature in Authorization header")
	}

	return credential, signature, nil
}

// createStringToSign creates string to sign in JavaScript script format
func createStringToSign(r *http.Request, body []byte, dateHeader, contentHash string) string {
	verb := r.Method
	pathAndQuery := r.URL.Path
	if r.URL.RawQuery != "" {
		pathAndQuery += "?" + r.URL.RawQuery
	}
	host := r.Host

	// Format: VERB\npathAndQuery\ntimestamp;host;contentHash
	stringToSign := verb + "\n" + pathAndQuery + "\n" + dateHeader + ";" + host + ";" + contentHash

	return stringToSign
}

// verifySignature verifies HMAC signature
func verifySignature(stringToSign, signature, secretKey string, algorithm func() hash.Hash) (bool, error) {
	// Decode secret from Base64
	secretBytes, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return false, fmt.Errorf("secret is not valid Base64: %v", err)
	}

	// Create HMAC
	h := hmac.New(algorithm, secretBytes)
	h.Write([]byte(stringToSign))
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Create logger for debug logging
	logger, err := logging.New(logging.LoadConfigFromEnv())
	if err != nil {
		logger = nil
	}

	if logger != nil {
		logger.Debug("HMAC signature verification",
			zap.Int("secret_length", len(secretBytes)),
			zap.String("string_to_sign", stringToSign),
			zap.String("expected_signature", expectedSignature),
			zap.String("received_signature", signature))
	}

	return hmac.Equal([]byte(signature), []byte(expectedSignature)), nil
}

// getAlgorithm returns hash function
func getAlgorithm(algorithm string) func() hash.Hash {
	switch algorithm {
	case "sha512":
		return sha512.New
	default:
		return sha256.New
	}
}

// findClientSecret finds client secret by its ID
func findClientSecret(clientID string, clientSecrets []HMACClientSecret) (string, error) {
	for _, client := range clientSecrets {
		if client.ClientID == clientID {
			return client.Secret, nil
		}
	}
	return "", fmt.Errorf("client ID not found: %s", clientID)
}

// checkRouteAccess checks client access rights to route
func checkRouteAccess(clientID, route string, routeRights HMACRouteRights) error {
	// If rights are not configured, allow access
	if len(routeRights) == 0 {
		return nil
	}

	// Get rights for client
	allowedRoutes, exists := routeRights[clientID]
	if !exists {
		return fmt.Errorf("client %s has no route permissions", clientID)
	}

	// Check route access
	for _, allowedRoute := range allowedRoutes {
		if allowedRoute == "*" || strings.HasPrefix(route, allowedRoute) {
			return nil
		}
	}

	return fmt.Errorf("client %s not authorized for route %s", clientID, route)
}

// CreateHMACSignature creates HMAC signature for request (for testing)
func CreateHMACSignature(method, path string, body []byte, dateHeader, contentHash, clientID, secretKey string, algorithm string) (string, error) {
	// Create string to sign
	stringToSign := method + "\n" + path + "\n" + dateHeader + ";localhost;" + contentHash

	// Get algorithm
	hashFunc := getAlgorithm(algorithm)

	// Decode secret from Base64 for HMAC creation
	secretBytes, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		// If not Base64 - this is an error!
		return "", fmt.Errorf("secret is not valid Base64: %v", err)
	}

	// Create signature
	h := hmac.New(hashFunc, secretBytes)
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// CreateSignedRequest creates HTTP request with HMAC signature (for testing)
func CreateSignedRequest(method, url string, body []byte, clientID, secretKey string, algorithm string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Create headers as in JavaScript script
	dateHeader := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	// Create content hash
	contentHash := createContentHash(body, algorithm)

	// Create signature
	signature, err := CreateHMACSignature(method, req.URL.Path, body, dateHeader, contentHash, clientID, secretKey, algorithm)
	if err != nil {
		return nil, err
	}

	// Create Authorization header
	auth := fmt.Sprintf("HMAC-SHA%s Credential=%s&SignedHeaders=date;host;x-content-hmac&Signature=%s",
		strings.TrimPrefix(algorithm, "sha"), clientID, signature)

	req.Header.Set("Date", dateHeader)
	req.Header.Set("x-content-hmac", contentHash)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// createContentHash creates content hash
func createContentHash(body []byte, algorithm string) string {
	var h hash.Hash
	switch algorithm {
	case "sha512":
		h = sha512.New()
	default:
		h = sha256.New()
	}
	h.Write(body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
