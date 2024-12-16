package merchant

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/model"
	"go.uber.org/zap"
)

func (b *business) ProcessCreateMerchant(ctx context.Context, request *api.CreateMerchantRequest) (error, uint32) {
	b.log.Info("ProcessCreateMerchant start")

	merchant := b.mappingCreateMerchant(request)

	err, merchantId := b.repository.CreateMerchant(ctx, merchant)

	b.log.Info("CreateMerchant done")

	if err != nil {
		b.log.Error("CreateMerchant error", zap.Error(err))
		return err, 0
	}

	b.log.Info("Merchant created", zap.Uint32("merchantAccountId", merchantId))

	b.log.Info("ProcessCreateMerchant end")

	return nil, merchantId
}

func (b *business) mappingCreateMerchant(request *api.CreateMerchantRequest) *model.Merchant {
	return &model.Merchant{
		Name:      request.GetName(),
		NetworkId: request.GetNetworkId(),
	}
}
