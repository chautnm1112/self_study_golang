package transaction

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/enum"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/model"
	"go.uber.org/zap"
)

func (b *business) ProcessCreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (error, uint32) {
	b.log.Info("ProcessCreateTransaction start")

	transaction := b.mappingCreateTransaction(request)

	err, transactionId := b.repository.CreateTransaction(ctx, nil, transaction)

	b.log.Info("CreateTransaction done")

	if err != nil {
		b.log.Error("CreateTransaction error", zap.Error(err))
		return err, 0
	}

	b.log.Info("ProcessCreateTransaction end")

	return nil, transactionId
}

func (b *business) mappingCreateTransaction(request *api.CreateTransactionRequest) *model.Transaction {
	return &model.Transaction{
		FromAccountId: request.GetFromAccountId(),
		ToAccountId:   request.GetToAccountId(),
		Points:        request.Point,
		Type:          mapTransactionType(request.Type),
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
