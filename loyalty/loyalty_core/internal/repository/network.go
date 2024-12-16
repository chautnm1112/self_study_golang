package repository

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/model"
	"gorm.io/gorm"
)

type NetworkRepository interface {
	CreateNetwork(ctx context.Context, network *model.Network) (error, uint32)
}

type networkRepo struct {
	*gorm.DB
	TableName string
}

func NewNetworkRepository(db *gorm.DB) NetworkRepository {
	return &networkRepo{
		db, model.Network{}.TableName(),
	}
}
func (r *networkRepo) CreateNetwork(ctx context.Context, network *model.Network) (error, uint32) {
	return r.WithContext(ctx).Table(r.TableName).Create(network).Error, network.ID
}
