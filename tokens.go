package americanexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// TokenService handles token management operations
type TokenService struct {
	client *Client
}

// NewTokenService creates a new token service
func NewTokenService(client *Client) *TokenService {
	return &TokenService{client: client}
}

// TokenRequest represents a token creation request
type TokenRequest struct {
	CardDetails  *CardDetails `json:"card_details"`
	CustomerID   string       `json:"customer_id,omitempty"`
	Description  string       `json:"description,omitempty"`
	SingleUse    bool         `json:"single_use,omitempty"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	ID           string    `json:"id"`
	Token        string    `json:"token"`
	CustomerID   string    `json:"customer_id"`
	Description  string    `json:"description"`
	CardLast4    string    `json:"card_last4"`
	CardBrand    string    `json:"card_brand"`
	ExpiryMonth  int       `json:"expiry_month"`
	ExpiryYear   int       `json:"expiry_year"`
	SingleUse    bool      `json:"single_use"`
	Used         bool      `json:"used"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// CreateToken creates a new payment token
func (ts *TokenService) CreateToken(ctx context.Context, req *TokenRequest) (*TokenResponse, error) {
	// Validate the token request
	if err := ValidateTokenRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	resp, err := ts.client.Post(ctx, "/tokens", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &token, nil
}

// GetToken retrieves a token by ID
func (ts *TokenService) GetToken(ctx context.Context, tokenID string) (*TokenResponse, error) {
	resp, err := ts.client.Get(ctx, fmt.Sprintf("/tokens/%s", tokenID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &token, nil
}

// DeleteToken deletes a token
func (ts *TokenService) DeleteToken(ctx context.Context, tokenID string) error {
	_, err := ts.client.Delete(ctx, fmt.Sprintf("/tokens/%s", tokenID))
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}

// ListTokensRequest represents parameters for listing tokens
type ListTokensRequest struct {
	CustomerID string `url:"customer_id,omitempty"`
	Limit      int    `url:"limit,omitempty"`
	Offset     int    `url:"offset,omitempty"`
}

// ListTokensResponse represents a list of tokens response
type ListTokensResponse struct {
	Tokens     []TokenResponse `json:"tokens"`
	TotalCount int             `json:"total_count"`
	HasMore    bool            `json:"has_more"`
}

// ListTokens retrieves a list of tokens
func (ts *TokenService) ListTokens(ctx context.Context, req *ListTokensRequest) (*ListTokensResponse, error) {
	query, err := encodeQuery(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	resp, err := ts.client.Get(ctx, "/tokens", query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokens ListTokensResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &tokens, nil
}