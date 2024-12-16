package repository

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *model.Transaction) (error, uint32)
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	CommitTransaction(ctx context.Context, tx *gorm.DB) error
	RollbackTransaction(ctx context.Context, tx *gorm.DB) error
}

type transactionRepo struct {
	*gorm.DB
	TableName string
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{
		db, model.Transaction{}.TableName(),
	}
}

func (r *transactionRepo) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (r *transactionRepo) CreateTransaction(ctx context.Context, tx *gorm.DB, transaction *model.Transaction) (error, uint32) {
	db := r.DB
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Table(r.TableName).Create(transaction).Error, transaction.ID
}

func (r *transactionRepo) CommitTransaction(ctx context.Context, tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepo) RollbackTransaction(ctx context.Context, tx *gorm.DB) error {
	if err := tx.Rollback().Error; err != nil {
		return err
	}
	return nil
}
