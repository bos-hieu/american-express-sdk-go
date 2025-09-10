package americanexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

// TransactionService handles transaction-related operations
type TransactionService struct {
	client *Client
}

// NewTransactionService creates a new transaction service
func NewTransactionService(client *Client) *TransactionService {
	return &TransactionService{client: client}
}

// TransactionRequest represents a transaction authorization request
type TransactionRequest struct {
	Amount       float64           `json:"amount"`
	Currency     string            `json:"currency"`
	MerchantID   string            `json:"merchant_id"`
	Description  string            `json:"description,omitempty"`
	Reference    string            `json:"reference,omitempty"`
	CardToken    string            `json:"card_token,omitempty"`
	CardDetails  *CardDetails      `json:"card_details,omitempty"`
	BillingAddr  *Address          `json:"billing_address,omitempty"`
	ShippingAddr *Address          `json:"shipping_address,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CaptureMode  string            `json:"capture_mode,omitempty"` // "auto", "manual"
	CVVCheck     bool              `json:"cvv_check,omitempty"`
	AVSCheck     bool              `json:"avs_check,omitempty"`
}

// TransactionResponse represents a transaction response
type TransactionResponse struct {
	ID                string            `json:"id"`
	Status            string            `json:"status"`
	Type              string            `json:"type"`
	Amount            float64           `json:"amount"`
	Currency          string            `json:"currency"`
	Description       string            `json:"description"`
	Reference         string            `json:"reference"`
	TransactionID     string            `json:"transaction_id"`
	AuthorizationCode string            `json:"authorization_code"`
	ProcessorResponse string            `json:"processor_response"`
	MerchantID        string            `json:"merchant_id"`
	CreatedAt         time.Time         `json:"created_at"`
	ProcessedAt       *time.Time        `json:"processed_at,omitempty"`
	ExpiresAt         *time.Time        `json:"expires_at,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	FailureReason     string            `json:"failure_reason,omitempty"`
	FailureCode       string            `json:"failure_code,omitempty"`
	CVVResult         string            `json:"cvv_result,omitempty"`
	AVSResult         string            `json:"avs_result,omitempty"`
}

// AuthorizeTransaction creates a new transaction authorization
func (ts *TransactionService) AuthorizeTransaction(ctx context.Context, req *TransactionRequest) (*TransactionResponse, error) {
	// Validate the transaction request
	if err := ValidateTransactionRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	resp, err := ts.client.Post(ctx, "/transactions/authorize", req)
	if err != nil {
		return nil, fmt.Errorf("failed to authorize transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var transaction TransactionResponse
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &transaction, nil
}

// GetTransaction retrieves a transaction by ID
func (ts *TransactionService) GetTransaction(ctx context.Context, transactionID string) (*TransactionResponse, error) {
	resp, err := ts.client.Get(ctx, fmt.Sprintf("/transactions/%s", transactionID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var transaction TransactionResponse
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &transaction, nil
}

// CaptureTransactionRequest represents a transaction capture request
type CaptureTransactionRequest struct {
	Amount    *float64          `json:"amount,omitempty"`
	Reference string            `json:"reference,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// CaptureTransaction captures a previously authorized transaction
func (ts *TransactionService) CaptureTransaction(ctx context.Context, transactionID string, req *CaptureTransactionRequest) (*TransactionResponse, error) {
	if req == nil {
		req = &CaptureTransactionRequest{}
	}

	resp, err := ts.client.Post(ctx, fmt.Sprintf("/transactions/%s/capture", transactionID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to capture transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var transaction TransactionResponse
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &transaction, nil
}

// VoidTransactionRequest represents a transaction void request
type VoidTransactionRequest struct {
	Reason    string            `json:"reason,omitempty"`
	Reference string            `json:"reference,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// VoidTransaction voids a previously authorized transaction
func (ts *TransactionService) VoidTransaction(ctx context.Context, transactionID string, req *VoidTransactionRequest) (*TransactionResponse, error) {
	if req == nil {
		req = &VoidTransactionRequest{}
	}

	resp, err := ts.client.Post(ctx, fmt.Sprintf("/transactions/%s/void", transactionID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to void transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var transaction TransactionResponse
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &transaction, nil
}

// RefundTransactionRequest represents a transaction refund request
type RefundTransactionRequest struct {
	Amount    float64           `json:"amount"`
	Reason    string            `json:"reason,omitempty"`
	Reference string            `json:"reference,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// RefundTransactionResponse represents a transaction refund response
type RefundTransactionResponse struct {
	ID                string            `json:"id"`
	TransactionID     string            `json:"transaction_id"`
	Amount            float64           `json:"amount"`
	Currency          string            `json:"currency"`
	Status            string            `json:"status"`
	Reason            string            `json:"reason"`
	Reference         string            `json:"reference"`
	RefundID          string            `json:"refund_id"`
	ProcessorResponse string            `json:"processor_response"`
	CreatedAt         time.Time         `json:"created_at"`
	ProcessedAt       *time.Time        `json:"processed_at,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	FailureReason     string            `json:"failure_reason,omitempty"`
	FailureCode       string            `json:"failure_code,omitempty"`
}

// RefundTransaction creates a refund for a transaction
func (ts *TransactionService) RefundTransaction(ctx context.Context, transactionID string, req *RefundTransactionRequest) (*RefundTransactionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("refund request is required")
	}

	resp, err := ts.client.Post(ctx, fmt.Sprintf("/transactions/%s/refund", transactionID), req)
	if err != nil {
		return nil, fmt.Errorf("failed to refund transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var refund RefundTransactionResponse
	if err := json.Unmarshal(body, &refund); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &refund, nil
}

// ListTransactionsRequest represents a request to list transactions
type ListTransactionsRequest struct {
	MerchantID  string `json:"merchant_id,omitempty"`
	Status      string `json:"status,omitempty"`
	Type        string `json:"type,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	Reference   string `json:"reference,omitempty"`
	MinAmount   string `json:"min_amount,omitempty"`
	MaxAmount   string `json:"max_amount,omitempty"`
	Currency    string `json:"currency,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Offset      int    `json:"offset,omitempty"`
	SortBy      string `json:"sort_by,omitempty"`
	SortOrder   string `json:"sort_order,omitempty"`
}

// ListTransactionsResponse represents a response with multiple transactions
type ListTransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int                   `json:"total"`
	Limit        int                   `json:"limit"`
	Offset       int                   `json:"offset"`
	HasMore      bool                  `json:"has_more"`
}

// ListTransactions retrieves a list of transactions with optional filters
func (ts *TransactionService) ListTransactions(ctx context.Context, req *ListTransactionsRequest) (*ListTransactionsResponse, error) {
	query := url.Values{}
	if req != nil {
		if req.MerchantID != "" {
			query.Add("merchant_id", req.MerchantID)
		}
		if req.Status != "" {
			query.Add("status", req.Status)
		}
		if req.Type != "" {
			query.Add("type", req.Type)
		}
		if req.StartDate != "" {
			query.Add("start_date", req.StartDate)
		}
		if req.EndDate != "" {
			query.Add("end_date", req.EndDate)
		}
		if req.Reference != "" {
			query.Add("reference", req.Reference)
		}
		if req.MinAmount != "" {
			query.Add("min_amount", req.MinAmount)
		}
		if req.MaxAmount != "" {
			query.Add("max_amount", req.MaxAmount)
		}
		if req.Currency != "" {
			query.Add("currency", req.Currency)
		}
		if req.Limit > 0 {
			query.Add("limit", fmt.Sprintf("%d", req.Limit))
		}
		if req.Offset > 0 {
			query.Add("offset", fmt.Sprintf("%d", req.Offset))
		}
		if req.SortBy != "" {
			query.Add("sort_by", req.SortBy)
		}
		if req.SortOrder != "" {
			query.Add("sort_order", req.SortOrder)
		}
	}

	resp, err := ts.client.Get(ctx, "/transactions", query)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var transactions ListTransactionsResponse
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &transactions, nil
}

// SearchTransactionsRequest represents a search request for transactions
type SearchTransactionsRequest struct {
	Query       string `json:"query"`
	MerchantID  string `json:"merchant_id,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Offset      int    `json:"offset,omitempty"`
}

// SearchTransactions searches for transactions using a query string
func (ts *TransactionService) SearchTransactions(ctx context.Context, req *SearchTransactionsRequest) (*ListTransactionsResponse, error) {
	if req == nil || req.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	query := url.Values{}
	query.Add("q", req.Query)
	if req.MerchantID != "" {
		query.Add("merchant_id", req.MerchantID)
	}
	if req.StartDate != "" {
		query.Add("start_date", req.StartDate)
	}
	if req.EndDate != "" {
		query.Add("end_date", req.EndDate)
	}
	if req.Limit > 0 {
		query.Add("limit", fmt.Sprintf("%d", req.Limit))
	}
	if req.Offset > 0 {
		query.Add("offset", fmt.Sprintf("%d", req.Offset))
	}

	resp, err := ts.client.Get(ctx, "/transactions/search", query)
	if err != nil {
		return nil, fmt.Errorf("failed to search transactions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var transactions ListTransactionsResponse
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &transactions, nil
}

// GetTransactionStatus retrieves the current status of a transaction
func (ts *TransactionService) GetTransactionStatus(ctx context.Context, transactionID string) (*TransactionResponse, error) {
	resp, err := ts.client.Get(ctx, fmt.Sprintf("/transactions/%s/status", transactionID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction status: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var transaction TransactionResponse
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &transaction, nil
}