package main

import (
	"context"
	"fmt"
	"log"

	amex "github.com/bos-hieu/american-express-sdk-go"
)

func main() {
	// Initialize the SDK
	config := &amex.Config{
		APIKey:    "your-api-key",
		SecretKey: "your-secret-key",
		// BaseURL: "https://gateway-na.americanexpress.com/api", // Optional, uses default
	}

	sdk := amex.NewSDK(config)

	ctx := context.Background()

	// Example 1: Create a payment token
	fmt.Println("Creating payment token...")
	tokenReq := &amex.TokenRequest{
		CardDetails: &amex.CardDetails{
			Number:      "4111111111111111",
			ExpiryMonth: 12,
			ExpiryYear:  2025,
			CVV:         "123",
			HolderName:  "John Doe",
		},
		CustomerID:  "customer_123",
		Description: "Payment token for John Doe",
		SingleUse:   false,
	}

	token, err := sdk.Tokens.CreateToken(ctx, tokenReq)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
	} else {
		fmt.Printf("Token created: %s (Last4: %s)\n", token.Token, token.CardLast4)
	}

	// Example 2: Create a payment using the token
	if token != nil {
		fmt.Println("\nCreating payment...")
		paymentReq := &amex.PaymentRequest{
			Amount:      100.00,
			Currency:    "USD",
			MerchantID:  "merchant_123",
			Description: "Test payment",
			Reference:   "order_456",
			CardToken:   token.Token,
			BillingAddr: &amex.Address{
				Line1:      "123 Main St",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "US",
			},
			Metadata: map[string]string{
				"order_id":     "456",
				"customer_id":  "123",
			},
		}

		payment, err := sdk.Payments.CreatePayment(ctx, paymentReq)
		if err != nil {
			log.Printf("Failed to create payment: %v", err)
		} else {
			fmt.Printf("Payment created: %s (Status: %s, Amount: %.2f %s)\n", 
				payment.ID, payment.Status, payment.Amount, payment.Currency)
		}

		// Example 3: Retrieve the payment
		if payment != nil {
			fmt.Println("\nRetrieving payment...")
			retrievedPayment, err := sdk.Payments.GetPayment(ctx, payment.ID)
			if err != nil {
				log.Printf("Failed to retrieve payment: %v", err)
			} else {
				fmt.Printf("Retrieved payment: %s (Status: %s)\n", 
					retrievedPayment.ID, retrievedPayment.Status)
			}
		}
	}

	// Example 4: Get merchant information
	fmt.Println("\nRetrieving merchant info...")
	merchantInfo, err := sdk.Merchant.GetMerchantInfo(ctx, "merchant_123")
	if err != nil {
		log.Printf("Failed to get merchant info: %v", err)
	} else {
		fmt.Printf("Merchant: %s (%s)\n", merchantInfo.Name, merchantInfo.Status)
	}

	// Example 5: List tokens
	fmt.Println("\nListing tokens...")
	listReq := &amex.ListTokensRequest{
		CustomerID: "customer_123",
		Limit:      10,
		Offset:     0,
	}

	tokenList, err := sdk.Tokens.ListTokens(ctx, listReq)
	if err != nil {
		log.Printf("Failed to list tokens: %v", err)
	} else {
		fmt.Printf("Found %d tokens (Total: %d)\n", len(tokenList.Tokens), tokenList.TotalCount)
		for _, t := range tokenList.Tokens {
			fmt.Printf("  Token %s: Last4 %s (%s)\n", t.ID, t.CardLast4, t.CardBrand)
		}
	}

	fmt.Println("\nExample completed!")
}