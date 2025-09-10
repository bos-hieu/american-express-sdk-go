package americanexpress

import (
	"testing"
)

func TestValidateCardDetails(t *testing.T) {
	tests := []struct {
		name    string
		card    *CardDetails
		wantErr bool
		errType error
	}{
		{
			name: "valid card",
			card: &CardDetails{
				Number:      "4111111111111111",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				CVV:         "123",
				HolderName:  "John Doe",
			},
			wantErr: false,
		},
		{
			name:    "nil card",
			card:    nil,
			wantErr: true,
		},
		{
			name: "invalid card number",
			card: &CardDetails{
				Number:      "123",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				CVV:         "123",
				HolderName:  "John Doe",
			},
			wantErr: true,
			errType: ErrInvalidCardNumber,
		},
		{
			name: "invalid expiry month",
			card: &CardDetails{
				Number:      "4111111111111111",
				ExpiryMonth: 13,
				ExpiryYear:  2025,
				CVV:         "123",
				HolderName:  "John Doe",
			},
			wantErr: true,
			errType: ErrInvalidExpiryDate,
		},
		{
			name: "invalid CVV",
			card: &CardDetails{
				Number:      "4111111111111111",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				CVV:         "12",
				HolderName:  "John Doe",
			},
			wantErr: true,
			errType: ErrInvalidCVV,
		},
		{
			name: "empty holder name",
			card: &CardDetails{
				Number:      "4111111111111111",
				ExpiryMonth: 12,
				ExpiryYear:  2025,
				CVV:         "123",
				HolderName:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCardDetails(tt.card)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCardDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestValidatePaymentRequest(t *testing.T) {
	validCard := &CardDetails{
		Number:      "4111111111111111",
		ExpiryMonth: 12,
		ExpiryYear:  2025,
		CVV:         "123",
		HolderName:  "John Doe",
	}

	tests := []struct {
		name    string
		req     *PaymentRequest
		wantErr bool
		errType error
	}{
		{
			name: "valid request with card details",
			req: &PaymentRequest{
				Amount:      100.00,
				Currency:    "USD",
				MerchantID:  "merchant_123",
				CardDetails: validCard,
			},
			wantErr: false,
		},
		{
			name: "valid request with token",
			req: &PaymentRequest{
				Amount:     100.00,
				Currency:   "USD",
				MerchantID: "merchant_123",
				CardToken:  "token_123",
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "invalid amount",
			req: &PaymentRequest{
				Amount:     0,
				Currency:   "USD",
				MerchantID: "merchant_123",
				CardToken:  "token_123",
			},
			wantErr: true,
			errType: ErrInvalidAmount,
		},
		{
			name: "invalid currency",
			req: &PaymentRequest{
				Amount:     100.00,
				Currency:   "US",
				MerchantID: "merchant_123",
				CardToken:  "token_123",
			},
			wantErr: true,
			errType: ErrInvalidCurrency,
		},
		{
			name: "missing payment method",
			req: &PaymentRequest{
				Amount:     100.00,
				Currency:   "USD",
				MerchantID: "merchant_123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePaymentRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePaymentRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestIsSupportedCurrency(t *testing.T) {
	tests := []struct {
		currency string
		want     bool
	}{
		{"USD", true},
		{"EUR", true},
		{"GBP", true},
		{"usd", true}, // case insensitive
		{"XYZ", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.currency, func(t *testing.T) {
			if got := IsSupportedCurrency(tt.currency); got != tt.want {
				t.Errorf("IsSupportedCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		want   float64
	}{
		{"whole number", 100.0, 100.0},
		{"two decimals", 100.25, 100.25},
		{"many decimals", 100.123456, 100.12},
		{"round up", 100.996, 100.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatAmount(tt.amount); got != tt.want {
				t.Errorf("FormatAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

