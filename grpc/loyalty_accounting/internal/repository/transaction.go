package repository

import (
	"context"
	"gorm.io/gorm"
	"loyalty_accounting/internal/model"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (error, uint32)
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

func (r *transactionRepo) CreateTransaction(ctx context.Context, transaction *model.Transaction) (error, uint32) {
	return r.WithContext(ctx).Table(r.TableName).Create(transaction).Error, transaction.ID
}
