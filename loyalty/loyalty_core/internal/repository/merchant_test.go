package repository

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMerchantRepository(t *testing.T) {
	merchantRepo := NewMerchantRepository(gormDb)

	t.Run("CreateMerchant", func(t *testing.T) {
		merchant := &model.Merchant{
			Name:      "Merchant 1",
			NetworkId: 1,
		}

		err, merchantID := merchantRepo.CreateMerchant(context.Background(), merchant)

		assert.NoError(t, err)
		assert.NotZero(t, merchantID)

		var createdMerchant model.Merchant
		err = gormDb.Table(createdMerchant.TableName()).Where("id = ?", merchantID).First(&createdMerchant).Error
		assert.NoError(t, err)
		assert.Equal(t, merchant.Name, createdMerchant.Name)
		assert.Equal(t, merchant.NetworkId, createdMerchant.NetworkId)
	})
}
