package enum

type AccountType string

const (
	TOTAL             AccountType = "TOTAL"
	NETWORK_PROMOTION AccountType = "NETWORK_PROMOTION"
	NETWORK_REVOKE    AccountType = "NETWORK_REVOKE"
)
