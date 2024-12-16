package account

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/repository"
	"go.uber.org/zap"
)

type Business interface {
	ProcessCreateAccount(ctx context.Context, request *api.CreateAccountRequest) (error, uint32)
	ProcessUpdateAccount(ctx context.Context, request *api.UpdateAccountRequest) (error, uint32)
	ProcessEarnPoints(ctx context.Context, request *api.EarnPointsRequest) (error, uint32)
	ProcessRedeemPoints(ctx context.Context, request *api.RedeemPointsRequest) (error, uint32)
	ProcessRefundEarnedPoints(ctx context.Context, request *api.RefundEarnedPointsRequest) (error, uint32)
	ProcessRefundRedeemPoints(ctx context.Context, request *api.RefundRedeemPointsRequest) (error, uint32)
}

type business struct {
	log        *zap.Logger
	repository *repository.Repository
}

func NewAccountBusiness(log *zap.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.Named("accountsBiz"),
		repository: repository,
	}
}
