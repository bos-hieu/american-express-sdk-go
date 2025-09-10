package americanexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// PaymentService handles payment-related operations
type PaymentService struct {
	client *Client
}

// NewPaymentService creates a new payment service
func NewPaymentService(client *Client) *PaymentService {
	return &PaymentService{client: client}
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	Amount       float64            `json:"amount"`
	Currency     string             `json:"currency"`
	MerchantID   string             `json:"merchant_id"`
	Description  string             `json:"description,omitempty"`
	Reference    string             `json:"reference,omitempty"`
	CardToken    string             `json:"card_token,omitempty"`
	CardDetails  *CardDetails       `json:"card_details,omitempty"`
	BillingAddr  *Address           `json:"billing_address,omitempty"`
	ShippingAddr *Address           `json:"shipping_address,omitempty"`
	Metadata     map[string]string  `json:"metadata,omitempty"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	ID                string            `json:"id"`
	Status            string            `json:"status"`
	Amount            float64           `json:"amount"`
	Currency          string            `json:"currency"`
	Description       string            `json:"description"`
	Reference         string            `json:"reference"`
	TransactionID     string            `json:"transaction_id"`
	AuthorizationCode string            `json:"authorization_code"`
	CreatedAt         time.Time         `json:"created_at"`
	ProcessedAt       *time.Time        `json:"processed_at,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	FailureReason     string            `json:"failure_reason,omitempty"`
}

// CardDetails represents card information
type CardDetails struct {
	Number      string `json:"number"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"`
	CVV         string `json:"cvv"`
	HolderName  string `json:"holder_name"`
}

// Address represents billing or shipping address
type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// CreatePayment creates a new payment
func (ps *PaymentService) CreatePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	// Validate the payment request
	if err := ValidatePaymentRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	resp, err := ps.client.Post(ctx, "/payments", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var payment PaymentResponse
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &payment, nil
}

// GetPayment retrieves a payment by ID
func (ps *PaymentService) GetPayment(ctx context.Context, paymentID string) (*PaymentResponse, error) {
	resp, err := ps.client.Get(ctx, fmt.Sprintf("/payments/%s", paymentID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var payment PaymentResponse
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &payment, nil
}

// CapturePayment captures an authorized payment
func (ps *PaymentService) CapturePayment(ctx context.Context, paymentID string, amount *float64) (*PaymentResponse, error) {
	captureReq := map[string]interface{}{}
	if amount != nil {
		captureReq["amount"] = *amount
	}

	resp, err := ps.client.Post(ctx, fmt.Sprintf("/payments/%s/capture", paymentID), captureReq)
	if err != nil {
		return nil, fmt.Errorf("failed to capture payment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var payment PaymentResponse
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &payment, nil
}

// VoidPayment voids an authorized payment
func (ps *PaymentService) VoidPayment(ctx context.Context, paymentID string) (*PaymentResponse, error) {
	resp, err := ps.client.Post(ctx, fmt.Sprintf("/payments/%s/void", paymentID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to void payment: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var payment PaymentResponse
	if err := json.Unmarshal(body, &payment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &payment, nil
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason,omitempty"`
	Reference string  `json:"reference,omitempty"`
}

// RefundResponse represents a refund response
type RefundResponse struct {
	ID          string    `json:"id"`
	PaymentID   string    `json:"payment_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	Reason      string    `json:"reason"`
	Reference   string    `json:"reference"`
	CreatedAt   time.Time `json:"created_at"`
	ProcessedAt time.Time `json:"processed_at"`
}

// CreateRefund creates a refund for a payment
func (ps *PaymentService) CreateRefund(ctx context.Context, req *RefundRequest) (*RefundResponse, error) {
	resp, err := ps.client.Post(ctx, "/refunds", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var refund RefundResponse
	if err := json.Unmarshal(body, &refund); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &refund, nil
}