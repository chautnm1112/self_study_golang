package merchant

import (
	"context"
	"loyalty_core/api"
	"loyalty_core/internal/model"
	"time"
)

func (b *business) ProcessCreateMerchant(ctx context.Context, request *api.CreateMerchantRequest) (error, uint32) {
	merchant := b.mappingCreateMerchant(request)

	return b.repository.CreateMerchant(ctx, merchant)
}

func (b *business) mappingCreateMerchant(request *api.CreateMerchantRequest) *model.Merchant {
	return &model.Merchant{
		Name:      request.GetName(),
		NetworkId: request.GetNetworkId(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
