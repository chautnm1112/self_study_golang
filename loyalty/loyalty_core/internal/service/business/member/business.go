package member

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/repository"
	"go.uber.org/zap"
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
