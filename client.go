// Package americanexpress provides a Go SDK for American Express APIs
package americanexpress

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for American Express APIs
	DefaultBaseURL = "https://gateway-na.americanexpress.com/api"
	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second
	// SDKVersion is the current version of this SDK
	SDKVersion = "1.0.0"
)

// Client represents the American Express API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	secretKey  string
	userAgent  string
}

// Config holds configuration for the American Express client
type Config struct {
	BaseURL    string
	APIKey     string
	SecretKey  string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// NewClient creates a new American Express API client
func NewClient(config *Config) *Client {
	if config == nil {
		config = &Config{}
	}

	// Set defaults
	if config.BaseURL == "" {
		config.BaseURL = DefaultBaseURL
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultTimeout
	}
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return &Client{
		baseURL:    strings.TrimSuffix(config.BaseURL, "/"),
		httpClient: config.HTTPClient,
		apiKey:     config.APIKey,
		secretKey:  config.SecretKey,
		userAgent:  fmt.Sprintf("AmexSDK-Go/%s", SDKVersion),
	}
}

// APIError represents an error response from the American Express API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Details    string `json:"details"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("amex api error: %d - %s (%s)", e.StatusCode, e.Message, e.Code)
}

// Request represents an HTTP request
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
	Query   url.Values
}

// doRequest executes an HTTP request and handles the response
func (c *Client) doRequest(ctx context.Context, req *Request) (*http.Response, error) {
	var body io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(jsonBody)
	}

	// Build URL
	reqURL := c.baseURL + req.Path
	if req.Query != nil && len(req.Query) > 0 {
		reqURL += "?" + req.Query.Encode()
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("User-Agent", c.userAgent)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Add authentication headers
	c.addAuthHeaders(httpReq)

	// Add custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for API errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		apiErr := &APIError{StatusCode: resp.StatusCode}
		
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			apiErr.Message = "failed to read error response"
		} else {
			// Try to parse error response
			if err := json.Unmarshal(respBody, apiErr); err != nil {
				apiErr.Message = string(respBody)
			}
		}
		
		return nil, apiErr
	}

	return resp, nil
}

// addAuthHeaders adds authentication headers to the request
func (c *Client) addAuthHeaders(req *http.Request) {
	if c.apiKey != "" {
		req.Header.Set("X-AMEX-API-KEY", c.apiKey)
	}
	// Additional authentication logic can be added here
	// This might include OAuth, JWT, or other authentication methods
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	return c.doRequest(ctx, &Request{
		Method: http.MethodGet,
		Path:   path,
		Query:  query,
	})
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.doRequest(ctx, &Request{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	})
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.doRequest(ctx, &Request{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	})
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.doRequest(ctx, &Request{
		Method: http.MethodDelete,
		Path:   path,
	})
}