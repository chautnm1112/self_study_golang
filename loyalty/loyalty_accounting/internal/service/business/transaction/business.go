package transaction

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/repository"
	"go.uber.org/zap"
)

type Business interface {
	ProcessCreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (error, uint32)
}

type business struct {
	log        *zap.Logger
	repository *repository.Repository
}

func NewTransaction(log *zap.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.Named("transactionBiz"),
		repository: repository,
	}
}
