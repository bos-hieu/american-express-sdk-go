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
    
    // Create a payment token
    tokenReq := &amex.TokenRequest{
        CardDetails: &amex.CardDetails{
            Number:      "4111111111111111",
            ExpiryMonth: 12,
            ExpiryYear:  2025,
            CVV:         "123",
            HolderName:  "John Doe",
        },
        CustomerID: "customer_123",
    }
    
    token, err := sdk.Tokens.CreateToken(ctx, tokenReq)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a payment
    paymentReq := &amex.PaymentRequest{
        Amount:     100.00,
        Currency:   "USD",
        MerchantID: "merchant_123",
        CardToken:  token.Token,
    }
    
    payment, err := sdk.Payments.CreatePayment(ctx, paymentReq)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Payment created: %s", payment.ID)
}
```

## Features

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

Check the `examples/` directory for more comprehensive examples:

- `examples/basic/main.go` - Basic usage examples
- More examples coming soon...

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