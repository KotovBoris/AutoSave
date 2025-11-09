package sbank

// API endpoints
const (
	endpointAuthToken           = "/auth/bank-token"
	endpointAccounts            = "/accounts"
	endpointAccountDetails      = "/accounts/%s"
	endpointAccountBalances     = "/accounts/%s/balances"
	endpointAccountTransactions = "/accounts/%s/transactions"
	endpointAccountStatus       = "/accounts/%s/status"
	endpointAccountClose        = "/accounts/%s/close"

	endpointAccountConsents = "/account-consents/request"
	endpointAccountConsent  = "/account-consents/%s"

	endpointPaymentConsents = "/payment-consents/request"
	endpointPaymentConsent  = "/payment-consents/%s"

	endpointPayments = "/payments"
	endpointPayment  = "/payments/%s"

	endpointProducts = "/products"
	endpointProduct  = "/products/%s"

	endpointProductAgreements = "/product-agreements"
	endpointProductAgreement  = "/product-agreements/%s"

	endpointProductAgreementConsents = "/product-agreement-consents/request"
	endpointProductAgreementConsent  = "/product-agreement-consents/%s"

	endpointCards      = "/cards"
	endpointCard       = "/cards/%s"
	endpointCardStatus = "/cards/%s/status"
	endpointCardLimits = "/cards/%s/limits"
)

