package americanexpress

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ErrInvalidCardNumber is returned when card number is invalid
	ErrInvalidCardNumber = errors.New("invalid card number")
	// ErrInvalidExpiryDate is returned when expiry date is invalid
	ErrInvalidExpiryDate = errors.New("invalid expiry date")
	// ErrInvalidCVV is returned when CVV is invalid
	ErrInvalidCVV = errors.New("invalid CVV")
	// ErrInvalidAmount is returned when amount is invalid
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrInvalidCurrency is returned when currency is invalid
	ErrInvalidCurrency = errors.New("invalid currency")
)

// cardNumberRegex matches basic card number patterns
var cardNumberRegex = regexp.MustCompile(`^\d{13,19}$`)

// ValidateCardDetails validates card details
func ValidateCardDetails(card *CardDetails) error {
	if card == nil {
		return errors.New("card details cannot be nil")
	}

	// Remove spaces and validate card number
	cardNumber := strings.ReplaceAll(card.Number, " ", "")
	if !cardNumberRegex.MatchString(cardNumber) {
		return ErrInvalidCardNumber
	}

	// Validate expiry date
	if card.ExpiryMonth < 1 || card.ExpiryMonth > 12 {
		return fmt.Errorf("%w: month must be 1-12", ErrInvalidExpiryDate)
	}
	if card.ExpiryYear < 2020 || card.ExpiryYear > 2099 {
		return fmt.Errorf("%w: year must be 2020-2099", ErrInvalidExpiryDate)
	}

	// Validate CVV
	if len(card.CVV) < 3 || len(card.CVV) > 4 {
		return ErrInvalidCVV
	}

	// Validate holder name
	if strings.TrimSpace(card.HolderName) == "" {
		return errors.New("holder name cannot be empty")
	}

	return nil
}

// ValidatePaymentRequest validates a payment request
func ValidatePaymentRequest(req *PaymentRequest) error {
	if req == nil {
		return errors.New("payment request cannot be nil")
	}

	// Validate amount
	if req.Amount <= 0 {
		return ErrInvalidAmount
	}

	// Validate currency
	if req.Currency == "" {
		return ErrInvalidCurrency
	}
	if len(req.Currency) != 3 {
		return fmt.Errorf("%w: currency must be 3 characters", ErrInvalidCurrency)
	}

	// Validate merchant ID
	if strings.TrimSpace(req.MerchantID) == "" {
		return errors.New("merchant ID cannot be empty")
	}

	// Validate that either card token or card details are provided
	if req.CardToken == "" && req.CardDetails == nil {
		return errors.New("either card token or card details must be provided")
	}

	// If card details are provided, validate them
	if req.CardDetails != nil {
		if err := ValidateCardDetails(req.CardDetails); err != nil {
			return fmt.Errorf("invalid card details: %w", err)
		}
	}

	return nil
}

// ValidateTokenRequest validates a token request
func ValidateTokenRequest(req *TokenRequest) error {
	if req == nil {
		return errors.New("token request cannot be nil")
	}

	if req.CardDetails == nil {
		return errors.New("card details are required for token creation")
	}

	return ValidateCardDetails(req.CardDetails)
}

// SupportedCurrencies returns a list of supported currencies
func SupportedCurrencies() []string {
	return []string{
		"USD", "EUR", "GBP", "CAD", "AUD", "JPY", "CHF", "SGD", "HKD", "SEK",
		"NOK", "DKK", "PLN", "CZK", "HUF", "ILS", "MXN", "BRL", "ARS", "CLP",
	}
}

// IsSupportedCurrency checks if a currency is supported
func IsSupportedCurrency(currency string) bool {
	supported := SupportedCurrencies()
	for _, c := range supported {
		if strings.EqualFold(c, currency) {
			return true
		}
	}
	return false
}

// FormatAmount formats an amount to 2 decimal places
func FormatAmount(amount float64) float64 {
	return float64(int(amount*100)) / 100
}