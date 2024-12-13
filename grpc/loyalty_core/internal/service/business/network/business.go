package network

import (
	"context"
	"go.uber.org/zap"
	"loyalty_core/api"
	"loyalty_core/internal/repository"
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
