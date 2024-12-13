package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"loyalty_accounting/api"
	"loyalty_accounting/internal/service/business"
)

type Service struct {
	log    *zap.Logger
	gormDb *gorm.DB
	biz    *business.Business
	api.UnimplementedLoyaltyAccountingServiceServer
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
