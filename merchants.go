package americanexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

// MerchantService handles merchant-related operations
type MerchantService struct {
	client *Client
}

// NewMerchantService creates a new merchant service
func NewMerchantService(client *Client) *MerchantService {
	return &MerchantService{client: client}
}

// MerchantInfo represents merchant information
type MerchantInfo struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Website       string    `json:"website"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       *Address  `json:"address"`
	BusinessType  string    `json:"business_type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetMerchantInfo retrieves merchant information
func (ms *MerchantService) GetMerchantInfo(ctx context.Context, merchantID string) (*MerchantInfo, error) {
	resp, err := ms.client.Get(ctx, fmt.Sprintf("/merchants/%s", merchantID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get merchant info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var merchant MerchantInfo
	if err := json.Unmarshal(body, &merchant); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &merchant, nil
}

// TransactionSummary represents transaction summary data
type TransactionSummary struct {
	Date            string  `json:"date"`
	TotalAmount     float64 `json:"total_amount"`
	TotalCount      int     `json:"total_count"`
	SuccessfulCount int     `json:"successful_count"`
	FailedCount     int     `json:"failed_count"`
	Currency        string  `json:"currency"`
}

// GetTransactionSummary retrieves transaction summary for a date range
func (ms *MerchantService) GetTransactionSummary(ctx context.Context, merchantID, startDate, endDate string) ([]TransactionSummary, error) {
	query := make(map[string]string)
	if startDate != "" {
		query["start_date"] = startDate
	}
	if endDate != "" {
		query["end_date"] = endDate
	}

	urlValues := url.Values{}
	for k, v := range query {
		urlValues.Add(k, v)
	}

	resp, err := ms.client.Get(ctx, fmt.Sprintf("/merchants/%s/transactions/summary", merchantID), urlValues)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction summary: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var summary []TransactionSummary
	if err := json.Unmarshal(body, &summary); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return summary, nil
}

// SettlementInfo represents settlement information
type SettlementInfo struct {
	ID          string    `json:"id"`
	MerchantID  string    `json:"merchant_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	SettledAt   time.Time `json:"settled_at"`
	CreatedAt   time.Time `json:"created_at"`
	Reference   string    `json:"reference"`
}

// GetSettlements retrieves settlement information
func (ms *MerchantService) GetSettlements(ctx context.Context, merchantID string, limit, offset int) ([]SettlementInfo, error) {
	query := url.Values{}
	if limit > 0 {
		query.Add("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		query.Add("offset", fmt.Sprintf("%d", offset))
	}

	resp, err := ms.client.Get(ctx, fmt.Sprintf("/merchants/%s/settlements", merchantID), query)
	if err != nil {
		return nil, fmt.Errorf("failed to get settlements: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var settlements []SettlementInfo
	if err := json.Unmarshal(body, &settlements); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return settlements, nil
}