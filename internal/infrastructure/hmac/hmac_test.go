package hmac

import (
	"testing"
)

func TestHMACService_New(t *testing.T) {
	secret := "test-secret"
	service := New(secret, SHA256)

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if string(service.secret) != secret {
		t.Errorf("Expected secret %s, got %s", secret, string(service.secret))
	}

	if service.algo != SHA256 {
		t.Errorf("Expected algorithm %s, got %s", SHA256, service.algo)
	}
}

func TestHMACService_NewWithSHA256(t *testing.T) {
	secret := "test-secret"
	service := NewWithSHA256(secret)

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if service.algo != SHA256 {
		t.Errorf("Expected algorithm %s, got %s", SHA256, service.algo)
	}
}

func TestHMACService_NewWithSHA512(t *testing.T) {
	secret := "test-secret"
	service := NewWithSHA512(secret)

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if service.algo != SHA512 {
		t.Errorf("Expected algorithm %s, got %s", SHA512, service.algo)
	}
}

func TestHMACService_Sign(t *testing.T) {
	secret := "test-secret"
	service := New(secret, SHA256)

	data := []byte("test-data")
	signature := service.Sign(data)

	if signature == "" {
		t.Fatal("Expected signature to be generated")
	}

	// Check that signature is always the same for the same data
	signature2 := service.Sign(data)
	if signature != signature2 {
		t.Error("Expected signatures to be identical for same data")
	}
}

func TestHMACService_SignString(t *testing.T) {
	secret := "test-secret"
	service := New(secret, SHA256)

	data := "test-data"
	signature := service.SignString(data)

	if signature == "" {
		t.Fatal("Expected signature to be generated")
	}

	// Check that string signature matches byte signature
	expectedSignature := service.Sign([]byte(data))
	if signature != expectedSignature {
		t.Error("Expected string signature to match byte signature")
	}
}

func TestHMACService_Verify(t *testing.T) {
	secret := "test-secret"
	service := New(secret, SHA256)

	data := []byte("test-data")
	signature := service.Sign(data)

	// Check correct signature
	if !service.Verify(data, signature) {
		t.Error("Expected valid signature to be verified")
	}

	// Check invalid signature
	if service.Verify(data, "invalid-signature") {
		t.Error("Expected invalid signature to be rejected")
	}

	// Check signature for other data
	otherData := []byte("other-data")
	if service.Verify(otherData, signature) {
		t.Error("Expected signature for different data to be rejected")
	}
}

func TestHMACService_VerifyString(t *testing.T) {
	secret := "test-secret"
	service := New(secret, SHA256)

	data := "test-data"
	signature := service.SignString(data)

	// Check correct signature
	if !service.VerifyString(data, signature) {
		t.Error("Expected valid signature to be verified")
	}

	// Check invalid signature
	if service.VerifyString(data, "invalid-signature") {
		t.Error("Expected invalid signature to be rejected")
	}
}

func TestHMACService_GetAlgorithm(t *testing.T) {
	secret := "test-secret"
	service := New(secret, SHA512)

	algo := service.GetAlgorithm()
	if algo != SHA512 {
		t.Errorf("Expected algorithm %s, got %s", SHA512, algo)
	}
}

func TestHMACService_DifferentAlgorithms(t *testing.T) {
	secret := "test-secret"
	data := "test-data"

	service256 := New(secret, SHA256)
	service512 := New(secret, SHA512)

	signature256 := service256.SignString(data)
	signature512 := service512.SignString(data)

	// Signatures should be different for different algorithms
	if signature256 == signature512 {
		t.Error("Expected different algorithms to produce different signatures")
	}

	// Check that each signature works with its own algorithm
	if !service256.VerifyString(data, signature256) {
		t.Error("Expected SHA256 signature to be verified by SHA256 service")
	}

	if !service512.VerifyString(data, signature512) {
		t.Error("Expected SHA512 signature to be verified by SHA512 service")
	}
}

func TestUtilityFunctions(t *testing.T) {
	secret := "test-secret"
	data := "test-data"

	// Test SignStringData
	signature, err := SignStringData(data, secret, SHA256)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if signature == "" {
		t.Fatal("Expected signature to be generated")
	}

	// Test VerifyStringData
	valid, err := VerifyStringData(data, signature, secret, SHA256)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !valid {
		t.Error("Expected signature to be verified")
	}

	// Test with empty secret
	_, err = SignStringData(data, "", SHA256)
	if err == nil {
		t.Error("Expected error for empty secret")
	}

	_, err = VerifyStringData(data, signature, "", SHA256)
	if err == nil {
		t.Error("Expected error for empty secret")
	}
}

func TestHMACService_Consistency(t *testing.T) {
	secret := "test-secret"
	service := New(secret, SHA256)

	data := "test-data"

	// Generate signature multiple times
	signatures := make([]string, 10)
	for i := 0; i < 10; i++ {
		signatures[i] = service.SignString(data)
	}

	// All signatures should be identical
	firstSignature := signatures[0]
	for i := 1; i < 10; i++ {
		if signatures[i] != firstSignature {
			t.Errorf("Signature %d differs from first signature", i)
		}
	}
}
