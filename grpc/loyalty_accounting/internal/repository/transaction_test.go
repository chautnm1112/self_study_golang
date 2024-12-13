package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	_ "gorm.io/gorm"
	"loyalty_accounting/internal/enum"
	"loyalty_accounting/internal/model"
	"testing"
)

func TestTransactionRepository(t *testing.T) {
	transactionRepo := NewTransactionRepository(gormDb)

	t.Run("CreateTransaction", func(t *testing.T) {
		transaction := &model.Transaction{
			FromAccountId: 1,
			ToAccountId:   2,
			Points:        1000,
			Type:          enum.EARN_POINTS,
		}

		err, transactionID := transactionRepo.CreateTransaction(context.Background(), transaction)

		assert.NoError(t, err)
		assert.NotZero(t, transactionID)

		var createdTransaction model.Transaction
		err = gormDb.Table(createdTransaction.TableName()).Where("id = ?", transactionID).First(&createdTransaction).Error
		assert.NoError(t, err)
		assert.Equal(t, transaction.FromAccountId, createdTransaction.FromAccountId)
		assert.Equal(t, transaction.ToAccountId, createdTransaction.ToAccountId)
		assert.Equal(t, transaction.Points, createdTransaction.Points)
		assert.Equal(t, transaction.Type, createdTransaction.Type)
	})
}
