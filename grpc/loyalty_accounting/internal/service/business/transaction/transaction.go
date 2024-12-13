package transaction

import (
	"context"
	"loyalty_accounting/api"
	"loyalty_accounting/internal/enum"
	"loyalty_accounting/internal/model"
	"time"
)

func (b *business) ProcessCreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (error, uint32) {
	transaction := b.mappingCreateTransaction(request)

	return b.repository.CreateTransaction(ctx, transaction)
}

func (b *business) mappingCreateTransaction(request *api.CreateTransactionRequest) *model.Transaction {
	return &model.Transaction{
		FromAccountId: request.GetFromAccountId(),
		ToAccountId:   request.GetToAccountId(),
		Points:        request.Point,
		Type:          mapTransactionType(request.Type),
		CreatedAt:     time.Now(),
	}
}

func mapTransactionType(protoType api.TransactionType) enum.TransactionType {
	switch protoType {
	case api.TransactionType_EARN_POINTS:
		return enum.EARN_POINTS
	case api.TransactionType_REDEEM_POINTS:
		return enum.REDEEM_POINTS
	case api.TransactionType_REFUND_EARNED_POINTS:
		return enum.REFUND_EARNED_POINTS
	case api.TransactionType_REFUND_REDEEMED_POINTS:
		return enum.REFUND_REDEEMED_POINTS
	default:
		return ""
	}
}
