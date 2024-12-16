package member

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/model"
	"go.uber.org/zap"
)

func (b *business) ProcessCreateMember(ctx context.Context, request *api.CreateMemberRequest) (error, uint32) {
	b.log.Info("ProcessCreateMember start")

	member := b.mappingCreateMember(request)

	err, memberId := b.repository.CreateMember(ctx, member)

	b.log.Info("CreateMember done")

	if err != nil {
		b.log.Error("CreateMember error", zap.Error(err))
		return err, 0
	}

	b.log.Info("Member created", zap.Uint32("memberAccountId", memberId))

	b.log.Info("ProcessCreateMember end")

	return nil, memberId
}

func (b *business) mappingCreateMember(request *api.CreateMemberRequest) *model.Member {
	return &model.Member{
		Name:      request.GetName(),
		Email:     request.GetEmail(),
		Phone:     request.GetPhone(),
		NetworkId: request.GetNetworkId(),
	}
}
