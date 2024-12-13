package repository

import "gorm.io/gorm"

type Repository struct {
	NetworkRepository
	MerchantRepository
	MemberRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		NewNetworkRepository(db),
		NewMerchantRepository(db),
		NewMemberRepository(db),
	}
}
