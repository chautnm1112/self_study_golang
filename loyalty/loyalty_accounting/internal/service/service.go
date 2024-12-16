package service

import (
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/service/business"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	log    *zap.Logger
	gormDb *gorm.DB
	biz    *business.Business
	api.UnimplementedLoyaltyAccountingServiceServer
}

func (s *Service) mustEmbedUnimplementedLoyaltyAccountingServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewService(
	logger *zap.Logger,
	gormDb *gorm.DB,
) (*Service, error) {
	biz := business.NewBusiness(logger, gormDb)
	return &Service{
		log:    logger,
		gormDb: gormDb,
		biz:    biz,
	}, nil
}
