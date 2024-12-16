package repository

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/enum"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestAccountRepository(t *testing.T) {
	ctx := context.Background()
	accountRepo := NewAccountRepository(gormDb)
	transactionRepo := NewTransactionRepository(gormDb)

	t.Run("CreateAccount", func(t *testing.T) {
		account := createTestAccount(1, enum.MEMBER, 100)

		tx, err := transactionRepo.BeginTransaction(ctx)
		assert.NoError(t, err, "Failed to begin transaction")

		defer rollbackOrCommit(ctx, transactionRepo, tx, &err, t)

		err, accountID := accountRepo.CreateAccount(ctx, tx, account)
		assert.NoError(t, err, "Failed to create account")

		createdAccount := &model.Account{}
		err = accountRepo.GetAccountById(ctx, tx, accountID, createdAccount)
		assert.NoError(t, err)
		assert.Equal(t, account.OwnerId, createdAccount.OwnerId)
		assert.Equal(t, account.Points, createdAccount.Points)
	})

	t.Run("GetAccountById", func(t *testing.T) {
		account := createTestAccount(6, enum.MERCHANT, 200)

		tx, err := transactionRepo.BeginTransaction(ctx)
		assert.NoError(t, err, "Failed to begin transaction")
		defer rollbackOrCommit(ctx, transactionRepo, tx, &err, t)

		err, accountID := accountRepo.CreateAccount(ctx, tx, account)
		assert.NoError(t, err, "Failed to create account")

		t.Log("accountID", accountID)

		foundAccount := &model.Account{}
		err = accountRepo.GetAccountById(ctx, tx, accountID, foundAccount)
		assert.NoError(t, err, "Failed to get account by ID")
		assert.Equal(t, account.OwnerId, foundAccount.OwnerId)
		assert.Equal(t, account.Points, foundAccount.Points)

		invalidAccount := &model.Account{}
		err = accountRepo.GetAccountById(ctx, tx, 9999, invalidAccount)
		assert.Error(t, err, "Expected error for non-existent account")
	})

	t.Run("UpdateAccount", func(t *testing.T) {
		account := createTestAccount(3, enum.MEMBER, 300)

		tx, err := transactionRepo.BeginTransaction(ctx)
		assert.NoError(t, err, "Failed to begin transaction")
		defer rollbackOrCommit(ctx, transactionRepo, tx, &err, t)

		err, accountID := accountRepo.CreateAccount(ctx, tx, account)
		assert.NoError(t, err, "Failed to create account")

		account.Points = 500
		err, _ = accountRepo.UpdateAccount(ctx, tx, account)
		assert.NoError(t, err, "Failed to update account")

		updatedAccount := &model.Account{}
		err = accountRepo.GetAccountById(ctx, tx, accountID, updatedAccount)
		assert.NoError(t, err)
		assert.Equal(t, uint32(500), updatedAccount.Points)
	})
}

func rollbackOrCommit(ctx context.Context, transactionRepo TransactionRepository, tx *gorm.DB, testErr *error, t *testing.T) {
	if *testErr != nil {
		rollbackErr := transactionRepo.RollbackTransaction(ctx, tx)
		assert.NoError(t, rollbackErr, "Failed to rollback transaction")
	} else {
		commitErr := transactionRepo.CommitTransaction(ctx, tx)
		assert.NoError(t, commitErr, "Failed to commit transaction")
	}
}

func createTestAccount(ownerId uint32, ownerType enum.OwnerType, points uint32) *model.Account {
	return &model.Account{
		OwnerId:   ownerId,
		OwnerType: ownerType,
		Points:    points,
	}
}
