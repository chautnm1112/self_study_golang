package transaction

import (
	"context"
	"go.uber.org/zap"
	"loyalty_accounting/api"
	"loyalty_accounting/internal/repository"
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
