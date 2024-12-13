package repository

import (
	"context"
	"gorm.io/gorm"
	"loyalty_core/internal/model"
)

type MerchantRepository interface {
	CreateMerchant(ctx context.Context, merchant *model.Merchant) (error, uint32)
}

type merchantRepo struct {
	*gorm.DB
	TableName string
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepo{
		db, model.Merchant{}.TableName(),
	}
}
func (r *merchantRepo) CreateMerchant(ctx context.Context, merchant *model.Merchant) (error, uint32) {
	return r.WithContext(ctx).Table(r.TableName).Create(merchant).Error, merchant.ID
}
