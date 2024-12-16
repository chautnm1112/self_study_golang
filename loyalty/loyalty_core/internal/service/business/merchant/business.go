package merchant

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/repository"
	"go.uber.org/zap"
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
