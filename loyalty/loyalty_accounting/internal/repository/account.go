package repository

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/model"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, tx *gorm.DB, account *model.Account) (error, uint32)
	GetAccountById(ctx context.Context, tx *gorm.DB, id uint32, pointsAccount *model.Account) error
	UpdateAccount(ctx context.Context, tx *gorm.DB, pointsAccount *model.Account) (error, uint32)
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

func (r *accountRepo) CreateAccount(ctx context.Context, tx *gorm.DB, account *model.Account) (error, uint32) {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Table(r.TableName).Create(account).Error, account.ID
}

func (r *accountRepo) GetAccountById(ctx context.Context, tx *gorm.DB, id uint32, account *model.Account) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Table(r.TableName).Where("id = ?", id).First(&account).Error
}

func (r *accountRepo) UpdateAccount(ctx context.Context, tx *gorm.DB, account *model.Account) (error, uint32) {
	db := r.DB
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Table(r.TableName).Where("id = ? AND deleted_at IS NULL", account.ID).Updates(account).Error, account.ID
}
