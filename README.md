# American Express SDK for Go

A comprehensive Go SDK for integrating with American Express APIs, providing easy-to-use interfaces for payment processing, tokenization, and merchant services.

## Installation

```bash
go get github.com/bos-hieu/american-express-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    amex "github.com/bos-hieu/american-express-sdk-go"
)

func main() {
    // Initialize the SDK
    config := &amex.Config{
        APIKey:    "your-api-key",
        SecretKey: "your-secret-key",
    }
    
    sdk := amex.NewSDK(config)
    ctx := context.Background()
    
    // Authorize a transaction
    transactionReq := &amex.TransactionRequest{
        Amount:      100.00,
        Currency:    "USD",
        MerchantID:  "merchant_123",
        Description: "Test transaction",
        CardDetails: &amex.CardDetails{
            Number:      "4111111111111111",
            ExpiryMonth: 12,
            ExpiryYear:  2025,
            CVV:         "123",
            HolderName:  "John Doe",
        },
        CaptureMode: "manual",
        CVVCheck:    true,
        AVSCheck:    true,
    }
    
    transaction, err := sdk.Transactions.AuthorizeTransaction(ctx, transactionReq)
    if err != nil {
        log.Fatal(err)
    }
    
    // Capture the transaction
    captureReq := &amex.CaptureTransactionRequest{
        Amount: &[]float64{50.00}[0], // Partial capture
    }
    
    captured, err := sdk.Transactions.CaptureTransaction(ctx, transaction.ID, captureReq)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Transaction captured: %s", captured.ID)
}
```

## Features

### Transaction Processing
- **Authorize transactions** with comprehensive validation and fraud checks
- **Capture transactions** (full or partial amounts)
- **Void authorized transactions** before capture
- **Refund transactions** with detailed tracking
- **List and search transactions** with flexible filtering
- **Get transaction status** and details
- **Advanced fraud protection** with CVV and AVS checks

### Payment Processing
- Create payments with card details or tokens
- Capture authorized payments
- Void payments
- Create refunds
- Retrieve payment information

### Token Management
- Create secure payment tokens
- Retrieve token information
- List customer tokens
- Delete tokens

### Merchant Services
- Retrieve merchant information
- Get transaction summaries
- Access settlement data

## Configuration

The SDK can be configured with various options:

```go
config := &amex.Config{
    APIKey:     "your-api-key",           // Required
    SecretKey:  "your-secret-key",        // Required
    BaseURL:    "custom-api-endpoint",    // Optional, defaults to production
    Timeout:    30 * time.Second,         // Optional, defaults to 30s
    HTTPClient: customHTTPClient,         // Optional, uses default client
}
```

## API Reference

### Transactions

#### Authorize Transaction
```go
transactionReq := &amex.TransactionRequest{
    Amount:      100.00,
    Currency:    "USD",
    MerchantID:  "merchant_123",
    Description: "Test transaction",
    Reference:   "order_123",
    CardDetails: &amex.CardDetails{
        Number:      "4111111111111111",
        ExpiryMonth: 12,
        ExpiryYear:  2025,
        CVV:         "123",
        HolderName:  "John Doe",
    },
    // Or use a card token instead
    // CardToken: "token_123",
    BillingAddr: &amex.Address{
        Line1:      "123 Main St",
        City:       "New York",
        State:      "NY",
        PostalCode: "10001",
        Country:    "US",
    },
    CaptureMode: "manual", // "auto" or "manual"
    CVVCheck:    true,
    AVSCheck:    true,
}

transaction, err := sdk.Transactions.AuthorizeTransaction(ctx, transactionReq)
```

#### Capture Transaction
```go
// Capture full amount
captured, err := sdk.Transactions.CaptureTransaction(ctx, transactionID, nil)

// Capture partial amount
captureReq := &amex.CaptureTransactionRequest{
    Amount:    &[]float64{50.00}[0],
    Reference: "partial_capture_123",
}
captured, err := sdk.Transactions.CaptureTransaction(ctx, transactionID, captureReq)
```

#### Void Transaction
```go
voidReq := &amex.VoidTransactionRequest{
    Reason:    "Customer cancelled order",
    Reference: "void_ref_123",
}
voided, err := sdk.Transactions.VoidTransaction(ctx, transactionID, voidReq)
```

#### Refund Transaction
```go
refundReq := &amex.RefundTransactionRequest{
    Amount:    25.00,
    Reason:    "Customer requested refund",
    Reference: "refund_ref_123",
}
refund, err := sdk.Transactions.RefundTransaction(ctx, transactionID, refundReq)
```

#### Get Transaction
```go
transaction, err := sdk.Transactions.GetTransaction(ctx, transactionID)

// Or get just the status
status, err := sdk.Transactions.GetTransactionStatus(ctx, transactionID)
```

#### List Transactions
```go
listReq := &amex.ListTransactionsRequest{
    MerchantID: "merchant_123",
    Status:     "authorized",
    StartDate:  "2023-01-01",
    EndDate:    "2023-01-31",
    Currency:   "USD",
    Limit:      10,
    Offset:     0,
    SortBy:     "created_at",
    SortOrder:  "desc",
}

transactions, err := sdk.Transactions.ListTransactions(ctx, listReq)
```

#### Search Transactions
```go
searchReq := &amex.SearchTransactionsRequest{
    Query:      "order_123",
    MerchantID: "merchant_123",
    StartDate:  "2023-01-01",
    EndDate:    "2023-01-31",
    Limit:      20,
}

results, err := sdk.Transactions.SearchTransactions(ctx, searchReq)
```

### Payments

#### Create Payment
```go
paymentReq := &amex.PaymentRequest{
    Amount:      100.00,
    Currency:    "USD",
    MerchantID:  "merchant_123",
    Description: "Test payment",
    CardToken:   "token_123",
    BillingAddr: &amex.Address{
        Line1:      "123 Main St",
        City:       "New York",
        State:      "NY",
        PostalCode: "10001",
        Country:    "US",
    },
}

payment, err := sdk.Payments.CreatePayment(ctx, paymentReq)
```

#### Capture Payment
```go
// Capture full amount
payment, err := sdk.Payments.CapturePayment(ctx, paymentID, nil)

// Capture partial amount
amount := 50.00
payment, err := sdk.Payments.CapturePayment(ctx, paymentID, &amount)
```

#### Create Refund
```go
refundReq := &amex.RefundRequest{
    PaymentID: "payment_123",
    Amount:    25.00,
    Reason:    "Customer requested refund",
}

refund, err := sdk.Payments.CreateRefund(ctx, refundReq)
```

### Tokens

#### Create Token
```go
tokenReq := &amex.TokenRequest{
    CardDetails: &amex.CardDetails{
        Number:      "4111111111111111",
        ExpiryMonth: 12,
        ExpiryYear:  2025,
        CVV:         "123",
        HolderName:  "John Doe",
    },
    CustomerID: "customer_123",
    SingleUse:  false,
}

token, err := sdk.Tokens.CreateToken(ctx, tokenReq)
```

#### List Tokens
```go
listReq := &amex.ListTokensRequest{
    CustomerID: "customer_123",
    Limit:      10,
    Offset:     0,
}

tokens, err := sdk.Tokens.ListTokens(ctx, listReq)
```

### Merchant Services

#### Get Merchant Info
```go
merchant, err := sdk.Merchant.GetMerchantInfo(ctx, "merchant_123")
```

#### Get Transaction Summary
```go
summary, err := sdk.Merchant.GetTransactionSummary(ctx, "merchant_123", "2023-01-01", "2023-01-31")
```

## Error Handling

The SDK provides structured error handling:

```go
payment, err := sdk.Payments.CreatePayment(ctx, paymentReq)
if err != nil {
    if apiErr, ok := err.(*amex.APIError); ok {
        log.Printf("API Error: %d - %s (%s)", apiErr.StatusCode, apiErr.Message, apiErr.Code)
    } else {
        log.Printf("Other error: %v", err)
    }
}
```

## Examples

Check the `examples/` directory for comprehensive examples:

- `examples/basic/main.go` - Basic usage examples
- `examples/transactions/main.go` - Complete transaction API examples including authorize, capture, void, refund, and search operations

## Testing

Run the tests:

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:
- Open an issue on GitHub
- Check the examples for common use cases
- Review the API documentation

## Version

Current version: v1.0.0

## Changelog

### v1.0.0
- Initial release
- Payment processing functionality
- Token management
- Merchant services
- Comprehensive error handling
- Examples and documentation