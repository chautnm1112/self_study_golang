package network

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/model"
	"go.uber.org/zap"
)

func (b *business) ProcessCreateNetwork(ctx context.Context, request *api.CreateNetworkRequest) (error, uint32) {
	b.log.Info("ProcessCreateNetwork start")

	network := b.mappingCreateNetwork(request)

	err, networkId := b.repository.CreateNetwork(ctx, network)

	b.log.Info("CreateNetwork done")

	if err != nil {
		b.log.Error("CreateNetwork error", zap.Error(err))
		return err, 0
	}

	b.log.Info("Network created", zap.Uint32("networkAccountId", networkId))

	b.log.Info("ProcessCreateNetwork end")

	return nil, networkId
}

func (b *business) mappingCreateNetwork(request *api.CreateNetworkRequest) *model.Network {
	return &model.Network{
		Name: request.GetName(),
	}
}
