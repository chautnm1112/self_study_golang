package member

import (
	"context"
	"go.uber.org/zap"
	"loyalty_core/api"
	"loyalty_core/internal/repository"
)

type Business interface {
	ProcessCreateMember(ctx context.Context, request *api.CreateMemberRequest) (error, uint32)
}

type business struct {
	log        *zap.Logger
	repository *repository.Repository
}

func NewMemberBusiness(log *zap.Logger, repository *repository.Repository) Business {
	return &business{
		log:        log.Named("memberBiz"),
		repository: repository,
	}
}
