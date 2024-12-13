package network

import (
	"context"
	"loyalty_core/api"
	"loyalty_core/internal/model"
	"time"
)

func (b *business) ProcessCreateNetwork(ctx context.Context, request *api.CreateNetworkRequest) (error, uint32) {
	network := b.mappingCreateNetwork(request)

	return b.repository.CreateNetwork(ctx, network)
}

func (b *business) mappingCreateNetwork(request *api.CreateNetworkRequest) *model.Network {
	return &model.Network{
		Name:      request.GetName(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
