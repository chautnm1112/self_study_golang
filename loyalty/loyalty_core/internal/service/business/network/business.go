package network

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/repository"
	"go.uber.org/zap"
)

type Business interface {
	ProcessCreateNetwork(ctx context.Context, request *api.CreateNetworkRequest) (error, uint32)
}

type business struct {
	log        *zap.Logger
	repository *repository.Repository
}

func NewNetWorkBusiness(log *zap.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.Named("networkBiz"),
		repository: repository,
	}
}
