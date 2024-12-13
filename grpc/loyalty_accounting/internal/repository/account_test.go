package repository

import (
	"context"
	_ "fmt"
	"github.com/stretchr/testify/assert"
	_ "github.com/testcontainers/testcontainers-go"
	_ "github.com/testcontainers/testcontainers-go/wait"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
	"loyalty_accounting/internal/enum"
	"loyalty_accounting/internal/model"
	"testing"
	_ "time"
)

func TestAccountRepository(t *testing.T) {
	ctx := context.Background()
	accountRepo := NewAccountRepository(gormDb)

	t.Run("CreateAccount", func(t *testing.T) {
		account := &model.Account{
			OwnerId:   1,
			OwnerType: enum.MEMBER,
			Points:    100,
		}
		err, accountID := accountRepo.CreateAccount(ctx, account)
		assert.NoError(t, err)
		assert.NotZero(t, accountID)

		var createdAccount model.Account
		err = gormDb.Table("accounts").Where("id = ?", accountID).First(&createdAccount).Error
		assert.NoError(t, err)
		assert.Equal(t, account.OwnerId, createdAccount.OwnerId)
		assert.Equal(t, account.OwnerType, createdAccount.OwnerType)
		assert.Equal(t, account.Points, createdAccount.Points)
	})

	t.Run("GetAccountByOwnerIdAndOwnerType", func(t *testing.T) {
		account := &model.Account{
			OwnerId:   2,
			OwnerType: enum.MERCHANT,
			Points:    200,
		}
		err, _ := accountRepo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		var foundAccount model.Account
		err = accountRepo.GetAccountByOwnerIdAndOwnerType(ctx, 2, enum.MERCHANT, &foundAccount)
		assert.NoError(t, err)
		assert.NotNil(t, foundAccount)
		assert.Equal(t, account.OwnerId, foundAccount.OwnerId)
	})

	t.Run("UpdateAccount", func(t *testing.T) {
		account := &model.Account{
			OwnerId:   3,
			OwnerType: enum.MEMBER,
			Points:    300,
		}
		err, accountID := accountRepo.CreateAccount(ctx, account)
		assert.NoError(t, err)

		account.Points = 500
		err, _ = accountRepo.UpdateAccount(ctx, account)
		assert.NoError(t, err)

		var updatedAccount model.Account
		err = gormDb.Table("accounts").Where("id = ?", accountID).First(&updatedAccount).Error
		assert.NoError(t, err)

		assert.Equal(t, int(500), int(updatedAccount.Points))
	})
}
