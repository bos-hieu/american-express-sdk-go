package americanexpress

// SDK represents the main American Express SDK client with all services
type SDK struct {
	*Client
	Payments     *PaymentService
	Tokens       *TokenService
	Merchant     *MerchantService
	Transactions *TransactionService
}

// NewSDK creates a new American Express SDK instance
func NewSDK(config *Config) *SDK {
	client := NewClient(config)
	
	return &SDK{
		Client:       client,
		Payments:     NewPaymentService(client),
		Tokens:       NewTokenService(client),
		Merchant:     NewMerchantService(client),
		Transactions: NewTransactionService(client),
	}
}

// Version returns the SDK version
func Version() string {
	return SDKVersion
}