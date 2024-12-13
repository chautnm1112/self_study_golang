package repository

import (
	"context"
	"gorm.io/gorm"
	"loyalty_accounting/internal/enum"
	"loyalty_accounting/internal/model"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *model.Account) (error, uint32)
	GetAccountByOwnerIdAndOwnerType(ctx context.Context, ownerId uint32, ownerType enum.OwnerType, pointsAccount *model.Account) error
	UpdateAccount(ctx context.Context, pointsAccount *model.Account) (error, uint32)
}

type accountRepo struct {
	*gorm.DB
	TableName string
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepo{
		db, model.Account{}.TableName(),
	}
}

func (r *accountRepo) CreateAccount(ctx context.Context, account *model.Account) (error, uint32) {
	return r.WithContext(ctx).Table(r.TableName).Create(account).Error, account.ID
}

func (r *accountRepo) GetAccountByOwnerIdAndOwnerType(ctx context.Context, ownerId uint32, ownerType enum.OwnerType, account *model.Account) error {
	return r.WithContext(ctx).Table(r.TableName).Where("owner_type = ? AND owner_id = ?", ownerType, ownerId).First(&account).Error
}

func (r *accountRepo) UpdateAccount(ctx context.Context, account *model.Account) (error, uint32) {
	return r.WithContext(ctx).Table(r.TableName).Updates(account).Error, account.ID
}
