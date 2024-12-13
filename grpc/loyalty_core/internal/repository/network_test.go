package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"loyalty_core/internal/model"
	"testing"
)

func TestNetworkRepository(t *testing.T) {
	networkRepo := NewNetworkRepository(gormDb)

	t.Run("CreateNetwork", func(t *testing.T) {
		network := &model.Network{
			Name: "Network One",
		}

		err, networkID := networkRepo.CreateNetwork(context.Background(), network)

		assert.NoError(t, err)
		assert.NotZero(t, networkID)

		var createdNetwork model.Network
		err = gormDb.Table(createdNetwork.TableName()).Where("id = ?", networkID).First(&createdNetwork).Error
		assert.NoError(t, err)
		assert.Equal(t, network.Name, createdNetwork.Name)
	})
}
