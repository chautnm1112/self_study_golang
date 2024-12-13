package merchant

import (
	"context"
	"go.uber.org/zap"
	"loyalty_core/api"
	"loyalty_core/internal/repository"
)

type Business interface {
	ProcessCreateMerchant(ctx context.Context, request *api.CreateMerchantRequest) (error, uint32)
}

type business struct {
	log        *zap.Logger
	repository *repository.Repository
}

func NewMerchantBusiness(log *zap.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.Named("merchantBiz"),
		repository: repository,
	}
}
