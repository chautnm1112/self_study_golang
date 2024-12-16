package repository

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/model"
	"github.com/stretchr/testify/assert"
	_ "gorm.io/gorm"
	"testing"
)

func TestMemberRepository(t *testing.T) {
	memberRepo := NewMemberRepository(gormDb)

	t.Run("CreateMember", func(t *testing.T) {
		member := &model.Member{
			Name:      "MinhChau",
			Email:     "chau.tnm@teko.vn",
			Phone:     "0346283293",
			NetworkId: 1,
		}

		err, memberID := memberRepo.CreateMember(context.Background(), member)

		assert.NoError(t, err)
		assert.NotZero(t, memberID)

		var createdMember model.Member
		err = gormDb.Table(createdMember.TableName()).Where("id = ?", memberID).First(&createdMember).Error
		assert.NoError(t, err)
		assert.Equal(t, member.Name, createdMember.Name)
		assert.Equal(t, member.Email, createdMember.Email)
		assert.Equal(t, member.Phone, createdMember.Phone)
		assert.Equal(t, member.NetworkId, createdMember.NetworkId)
	})
}
