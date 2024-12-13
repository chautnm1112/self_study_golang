package business

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"loyalty_core/internal/repository"
	"loyalty_core/internal/service/business/member"
	"loyalty_core/internal/service/business/merchant"
	"loyalty_core/internal/service/business/network"
)

type Business struct {
	db               *gorm.DB
	repository       *repository.Repository
	NetworkBusiness  network.Business
	MerchantBusiness merchant.Business
	MemberBusiness   member.Business
}

func NewBusiness(
	log *zap.Logger,
	db *gorm.DB,
) *Business {
	repo := repository.NewRepository(db)
	return &Business{
		db:               db,
		repository:       repo,
		NetworkBusiness:  network.NewNetWorkBusiness(log, repo),
		MerchantBusiness: merchant.NewMerchantBusiness(log, repo),
		MemberBusiness:   member.NewMemberBusiness(log, repo),
	}
}
