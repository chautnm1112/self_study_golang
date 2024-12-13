package member

import (
	"context"
	"loyalty_core/api"
	"loyalty_core/internal/model"
	"time"
)

func (b *business) ProcessCreateMember(ctx context.Context, request *api.CreateMemberRequest) (error, uint32) {
	b.log.Info("ProcessCreateMember")
	member := b.mappingCreateMember(request)

	return b.repository.CreateMember(ctx, member)
}

func (b *business) mappingCreateMember(request *api.CreateMemberRequest) *model.Member {
	return &model.Member{
		Name:      request.GetName(),
		Email:     request.GetEmail(),
		Phone:     request.GetPhone(),
		NetworkId: request.GetNetworkId(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
