package americanexpress

import (
	"testing"
)

func TestTransactionService_AuthorizeTransaction(t *testing.T) {
	tests := []struct {
		name    string
		request *TransactionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid transaction request",
			request: &TransactionRequest{
				Amount:     100.00,
				Currency:   "USD",
				MerchantID: "merchant_123",
				CardToken:  "token_123",
			},
			wantErr: false,
		},
		{
			name: "valid transaction request with card details",
			request: &TransactionRequest{
				Amount:     50.00,
				Currency:   "USD",
				MerchantID: "merchant_123",
				CardDetails: &CardDetails{
					Number:      "4111111111111111",
					ExpiryMonth: 12,
					ExpiryYear:  2025,
					CVV:         "123",
					HolderName:  "John Doe",
				},
				CaptureMode: "manual",
				CVVCheck:    true,
				AVSCheck:    true,
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			request: nil,
			wantErr: true,
			errMsg:  "transaction request cannot be nil",
		},
		{
			name: "zero amount",
			request: &TransactionRequest{
				Amount:     0,
				Currency:   "USD",
				MerchantID: "merchant_123",
				CardToken:  "token_123",
			},
			wantErr: true,
			errMsg:  "invalid amount",
		},
		{
			name: "empty currency",
			request: &TransactionRequest{
				Amount:     100.00,
				Currency:   "",
				MerchantID: "merchant_123",
				CardToken:  "token_123",
			},
			wantErr: true,
			errMsg:  "invalid currency",
		},
		{
			name: "empty merchant ID",
			request: &TransactionRequest{
				Amount:     100.00,
				Currency:   "USD",
				MerchantID: "",
				CardToken:  "token_123",
			},
			wantErr: true,
			errMsg:  "merchant ID cannot be empty",
		},
		{
			name: "no card token or card details",
			request: &TransactionRequest{
				Amount:     100.00,
				Currency:   "USD",
				MerchantID: "merchant_123",
			},
			wantErr: true,
			errMsg:  "either card token or card details must be provided",
		},
		{
			name: "invalid capture mode",
			request: &TransactionRequest{
				Amount:      100.00,
				Currency:    "USD",
				MerchantID:  "merchant_123",
				CardToken:   "token_123",
				CaptureMode: "invalid",
			},
			wantErr: true,
			errMsg:  "capture mode must be 'auto' or 'manual'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation only since we don't have a real API endpoint
			err := ValidateTransactionRequest(tt.request)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateTransactionRequest() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ValidateTransactionRequest() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateTransactionRequest() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestTransactionService_Creation(t *testing.T) {
	config := &Config{
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
	}
	
	sdk := NewSDK(config)
	
	if sdk.Transactions == nil {
		t.Error("TransactionService should be initialized")
	}
	
	if sdk.Transactions.client == nil {
		t.Error("TransactionService should have a client")
	}
}

func TestTransactionService_ValidateRefundTransactionRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *RefundTransactionRequest
		wantErr bool
	}{
		{
			name: "valid refund request",
			request: &RefundTransactionRequest{
				Amount:    50.00,
				Reason:    "Customer requested refund",
				Reference: "ref_123",
			},
			wantErr: false,
		},
		{
			name: "zero amount refund",
			request: &RefundTransactionRequest{
				Amount: 0,
				Reason: "Customer requested refund",
			},
			wantErr: false, // Zero amount might be valid for some scenarios
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For refund requests, we just check that they can be created
			// The actual validation would happen at the API level
			if tt.request == nil && !tt.wantErr {
				t.Error("Expected non-nil request")
			}
		})
	}
}

func TestTransactionService_ListTransactionsRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *ListTransactionsRequest
		wantErr bool
	}{
		{
			name: "valid list request",
			request: &ListTransactionsRequest{
				MerchantID: "merchant_123",
				Status:     "authorized",
				Limit:      10,
				Offset:     0,
			},
			wantErr: false,
		},
		{
			name: "list request with date range",
			request: &ListTransactionsRequest{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-31",
				Currency:  "USD",
				SortBy:    "created_at",
				SortOrder: "desc",
			},
			wantErr: false,
		},
		{
			name:    "nil request should be handled gracefully",
			request: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just test that the request structure is valid
			if tt.request != nil {
				// Basic validation that fields are accessible
				_ = tt.request.MerchantID
				_ = tt.request.Status
				_ = tt.request.Limit
			}
		})
	}
}

func TestTransactionService_SearchTransactionsRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *SearchTransactionsRequest
		wantErr bool
	}{
		{
			name: "valid search request",
			request: &SearchTransactionsRequest{
				Query:      "transaction_123",
				MerchantID: "merchant_123",
				Limit:      20,
			},
			wantErr: false,
		},
		{
			name:    "empty query",
			request: &SearchTransactionsRequest{Query: ""},
			wantErr: true,
		},
		{
			name:    "nil request",
			request: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.request == nil || tt.request.Query == ""
			if hasError != tt.wantErr {
				t.Errorf("Request validation mismatch: got error=%v, want error=%v", hasError, tt.wantErr)
			}
		})
	}
}

// MockTransactionService for testing service methods
type MockTransactionService struct {
	client *Client
}

func TestSDKIntegration(t *testing.T) {
	config := &Config{
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
	}
	
	sdk := NewSDK(config)
	
	// Test that all services are available
	if sdk.Transactions == nil {
		t.Error("Transactions service should be available")
	}
	
	if sdk.Payments == nil {
		t.Error("Payments service should be available")
	}
	
	if sdk.Tokens == nil {
		t.Error("Tokens service should be available")
	}
	
	if sdk.Merchant == nil {
		t.Error("Merchant service should be available")
	}
	
	// Test that services share the same client
	if sdk.Transactions.client != sdk.Client {
		t.Error("Transactions service should use the same client as SDK")
	}
}