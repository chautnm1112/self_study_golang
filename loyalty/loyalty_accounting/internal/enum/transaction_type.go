package enum

type TransactionType string

const (
	EARN_POINTS            TransactionType = "EARN_POINTS"
	REDEEM_POINTS          TransactionType = "REDEEM_POINTS"
	REFUND_EARNED_POINTS   TransactionType = "REFUND_EARNED_POINTS"
	REFUND_REDEEMED_POINTS TransactionType = "REFUND_REDEEMED_POINTS"
)
