package americanexpress

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := &Config{
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
		Timeout:   10 * time.Second,
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}

	if client.apiKey != "test-api-key" {
		t.Errorf("Expected API key to be 'test-api-key', got '%s'", client.apiKey)
	}

	if client.secretKey != "test-secret-key" {
		t.Errorf("Expected secret key to be 'test-secret-key', got '%s'", client.secretKey)
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("Expected base URL to be '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}
}

func TestNewClientWithDefaults(t *testing.T) {
	client := NewClient(nil)

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}

	if client.baseURL != DefaultBaseURL {
		t.Errorf("Expected base URL to be '%s', got '%s'", DefaultBaseURL, client.baseURL)
	}

	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("Expected timeout to be %v, got %v", DefaultTimeout, client.httpClient.Timeout)
	}
}

func TestNewSDK(t *testing.T) {
	config := &Config{
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
	}

	sdk := NewSDK(config)

	if sdk == nil {
		t.Fatal("Expected SDK to be non-nil")
	}

	if sdk.Client == nil {
		t.Fatal("Expected SDK client to be non-nil")
	}

	if sdk.Payments == nil {
		t.Fatal("Expected payments service to be non-nil")
	}

	if sdk.Tokens == nil {
		t.Fatal("Expected tokens service to be non-nil")
	}

	if sdk.Merchant == nil {
		t.Fatal("Expected merchant service to be non-nil")
	}
}

func TestVersion(t *testing.T) {
	version := Version()
	if version != SDKVersion {
		t.Errorf("Expected version to be '%s', got '%s'", SDKVersion, version)
	}
}

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: 400,
		Message:    "Bad Request",
		Code:       "INVALID_REQUEST",
		Details:    "Missing required field",
	}

	expected := "amex api error: 400 - Bad Request (INVALID_REQUEST)"
	if err.Error() != expected {
		t.Errorf("Expected error message to be '%s', got '%s'", expected, err.Error())
	}
}