package repository

import (
	"context"
	"gorm.io/gorm"
	"loyalty_core/internal/model"
)

type MemberRepository interface {
	CreateMember(ctx context.Context, member *model.Member) (error, uint32)
}

type memberRepo struct {
	*gorm.DB
	TableName string
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepo{
		db, model.Member{}.TableName(),
	}
}

func (r *memberRepo) CreateMember(ctx context.Context, member *model.Member) (error, uint32) {
	return r.WithContext(ctx).Table(r.TableName).Create(member).Error, member.ID
}
