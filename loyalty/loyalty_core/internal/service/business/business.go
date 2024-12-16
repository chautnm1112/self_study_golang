package business

import (
	"github.com/chautnm1112/loyalty/loyalty_core/internal/repository"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/service/business/member"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/service/business/merchant"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/service/business/network"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
