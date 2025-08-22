package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

// Algorithm represents supported hashing algorithms
type Algorithm string

const (
	SHA256 Algorithm = "sha256"
	SHA512 Algorithm = "sha512"
)

// HMACService provides methods for working with HMAC signatures
type HMACService struct {
	secret []byte
	algo   Algorithm
}

// New creates a new HMACService instance
func New(secret string, algorithm Algorithm) *HMACService {
	return &HMACService{
		secret: []byte(secret),
		algo:   algorithm,
	}
}

// NewWithSHA256 creates HMACService with SHA256 algorithm
func NewWithSHA256(secret string) *HMACService {
	return New(secret, SHA256)
}

// NewWithSHA512 creates HMACService with SHA512 algorithm
func NewWithSHA512(secret string) *HMACService {
	return New(secret, SHA512)
}

// Sign creates HMAC signature for data
func (h *HMACService) Sign(data []byte) string {
	hash := h.getHash()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// SignString creates HMAC signature for string
func (h *HMACService) SignString(data string) string {
	return h.Sign([]byte(data))
}

// Verify verifies HMAC signature
func (h *HMACService) Verify(data []byte, signature string) bool {
	expectedSignature := h.Sign(data)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// VerifyString verifies HMAC signature for string
func (h *HMACService) VerifyString(data string, signature string) bool {
	return h.Verify([]byte(data), signature)
}

// GetAlgorithm returns the used algorithm
func (h *HMACService) GetAlgorithm() Algorithm {
	return h.algo
}

// getHash returns the corresponding hash function
func (h *HMACService) getHash() hash.Hash {
	switch h.algo {
	case SHA256:
		return hmac.New(sha256.New, h.secret)
	case SHA512:
		return hmac.New(sha512.New, h.secret)
	default:
		// Use SHA256 by default
		return hmac.New(sha256.New, h.secret)
	}
}

// Utility functions for quick use without creating a service

// SignData creates HMAC signature for data with specified secret and algorithm
func SignData(data []byte, secret string, algorithm Algorithm) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("secret cannot be empty")
	}

	service := New(secret, algorithm)
	return service.Sign(data), nil
}

// SignStringData creates HMAC signature for string with specified secret and algorithm
func SignStringData(data string, secret string, algorithm Algorithm) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("secret cannot be empty")
	}

	service := New(secret, algorithm)
	return service.SignString(data), nil
}

// VerifyData verifies HMAC signature for data
func VerifyData(data []byte, signature string, secret string, algorithm Algorithm) (bool, error) {
	if secret == "" {
		return false, fmt.Errorf("secret cannot be empty")
	}

	service := New(secret, algorithm)
	return service.Verify(data, signature), nil
}

// VerifyStringData verifies HMAC signature for string
func VerifyStringData(data string, signature string, secret string, algorithm Algorithm) (bool, error) {
	if secret == "" {
		return false, fmt.Errorf("secret cannot be empty")
	}

	service := New(secret, algorithm)
	return service.VerifyString(data, signature), nil
}
