package main

import (
	"context"
	"log"
	"time"

	amex "github.com/bos-hieu/american-express-sdk-go"
)

func main() {
	// Initialize the SDK
	config := &amex.Config{
		APIKey:    "your-api-key",
		SecretKey: "your-secret-key",
		// BaseURL: "https://sandbox-gateway-na.americanexpress.com/api", // For sandbox
	}

	sdk := amex.NewSDK(config)
	ctx := context.Background()

	// Example 1: Authorize a transaction
	log.Println("=== Authorizing Transaction ===")
	transactionReq := &amex.TransactionRequest{
		Amount:      100.00,
		Currency:    "USD",
		MerchantID:  "merchant_123",
		Description: "Test purchase",
		Reference:   "order_12345",
		CardDetails: &amex.CardDetails{
			Number:      "4111111111111111",
			ExpiryMonth: 12,
			ExpiryYear:  2025,
			CVV:         "123",
			HolderName:  "John Doe",
		},
		BillingAddr: &amex.Address{
			Line1:      "123 Main Street",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "US",
		},
		CaptureMode: "manual", // We'll capture manually later
		CVVCheck:    true,
		AVSCheck:    true,
	}

	transaction, err := sdk.Transactions.AuthorizeTransaction(ctx, transactionReq)
	if err != nil {
		log.Printf("Failed to authorize transaction: %v", err)
		return
	}

	log.Printf("Transaction authorized: %s", transaction.ID)
	log.Printf("Authorization Code: %s", transaction.AuthorizationCode)
	log.Printf("Status: %s", transaction.Status)

	// Example 2: Get transaction details
	log.Println("\n=== Getting Transaction Details ===")
	retrievedTransaction, err := sdk.Transactions.GetTransaction(ctx, transaction.ID)
	if err != nil {
		log.Printf("Failed to get transaction: %v", err)
		return
	}

	log.Printf("Retrieved transaction: %s", retrievedTransaction.ID)
	log.Printf("Amount: %.2f %s", retrievedTransaction.Amount, retrievedTransaction.Currency)

	// Example 3: Capture the transaction (partial capture)
	log.Println("\n=== Capturing Transaction (Partial) ===")
	partialAmount := 75.00
	captureReq := &amex.CaptureTransactionRequest{
		Amount:    &partialAmount,
		Reference: "partial_capture_12345",
		Metadata: map[string]string{
			"capture_type": "partial",
			"reason":       "partial_shipment",
		},
	}

	captured, err := sdk.Transactions.CaptureTransaction(ctx, transaction.ID, captureReq)
	if err != nil {
		log.Printf("Failed to capture transaction: %v", err)
		return
	}

	log.Printf("Transaction captured: %s", captured.ID)
	log.Printf("Captured amount: %.2f", *captureReq.Amount)

	// Example 4: Create a refund
	log.Println("\n=== Creating Refund ===")
	refundReq := &amex.RefundTransactionRequest{
		Amount:    25.00,
		Reason:    "Customer returned item",
		Reference: "refund_12345",
		Metadata: map[string]string{
			"return_reason": "defective_item",
			"return_date":   time.Now().Format("2006-01-02"),
		},
	}

	refund, err := sdk.Transactions.RefundTransaction(ctx, transaction.ID, refundReq)
	if err != nil {
		log.Printf("Failed to create refund: %v", err)
		return
	}

	log.Printf("Refund created: %s", refund.ID)
	log.Printf("Refund amount: %.2f", refund.Amount)

	// Example 5: List transactions with filters
	log.Println("\n=== Listing Transactions ===")
	listReq := &amex.ListTransactionsRequest{
		MerchantID: "merchant_123",
		Status:     "captured",
		Currency:   "USD",
		StartDate:  "2023-01-01",
		EndDate:    time.Now().Format("2006-01-02"),
		Limit:      10,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}

	transactions, err := sdk.Transactions.ListTransactions(ctx, listReq)
	if err != nil {
		log.Printf("Failed to list transactions: %v", err)
		return
	}

	log.Printf("Found %d transactions", len(transactions.Transactions))
	for i, tx := range transactions.Transactions {
		if i >= 3 { // Limit output for example
			break
		}
		log.Printf("  %d. ID: %s, Amount: %.2f, Status: %s", 
			i+1, tx.ID, tx.Amount, tx.Status)
	}

	// Example 6: Search transactions
	log.Println("\n=== Searching Transactions ===")
	searchReq := &amex.SearchTransactionsRequest{
		Query:      "order_12345",
		MerchantID: "merchant_123",
		Limit:      5,
	}

	searchResults, err := sdk.Transactions.SearchTransactions(ctx, searchReq)
	if err != nil {
		log.Printf("Failed to search transactions: %v", err)
		return
	}

	log.Printf("Search found %d transactions", len(searchResults.Transactions))

	// Example 7: Authorize with token instead of card details
	log.Println("\n=== Authorizing with Token ===")
	
	// First create a token
	tokenReq := &amex.TokenRequest{
		CardDetails: &amex.CardDetails{
			Number:      "4111111111111111",
			ExpiryMonth: 11,
			ExpiryYear:  2026,
			CVV:         "456",
			HolderName:  "Jane Smith",
		},
		CustomerID: "customer_456",
		SingleUse:  false,
	}

	token, err := sdk.Tokens.CreateToken(ctx, tokenReq)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		return
	}

	// Now authorize using the token
	tokenTransactionReq := &amex.TransactionRequest{
		Amount:      200.00,
		Currency:    "USD",
		MerchantID:  "merchant_123",
		Description: "Token-based transaction",
		Reference:   "token_order_67890",
		CardToken:   token.Token,
		CaptureMode: "auto", // Auto-capture this time
	}

	tokenTransaction, err := sdk.Transactions.AuthorizeTransaction(ctx, tokenTransactionReq)
	if err != nil {
		log.Printf("Failed to authorize token transaction: %v", err)
		return
	}

	log.Printf("Token transaction authorized: %s", tokenTransaction.ID)
	log.Printf("Auto-captured: %s", tokenTransaction.Status)

	// Example 8: Check transaction status
	log.Println("\n=== Checking Transaction Status ===")
	status, err := sdk.Transactions.GetTransactionStatus(ctx, tokenTransaction.ID)
	if err != nil {
		log.Printf("Failed to get transaction status: %v", err)
		return
	}

	log.Printf("Transaction %s status: %s", status.ID, status.Status)
	if status.ProcessedAt != nil {
		log.Printf("Processed at: %s", status.ProcessedAt.Format(time.RFC3339))
	}

	log.Println("\n=== All Examples Completed ===")
}