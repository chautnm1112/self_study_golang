package service

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"loyalty_core/api"
	"loyalty_core/internal/service/business"
)

type Service struct {
	log    *zap.Logger
	gormDb *gorm.DB
	biz    *business.Business
	api.UnsafeLoyaltyCoreServiceServer
	api.UnimplementedLoyaltyAccountingServiceServer
	client api.LoyaltyAccountingServiceClient
}

func NewService(
	logger *zap.Logger,
	gormDb *gorm.DB,
	client api.LoyaltyAccountingServiceClient,
) (*Service, error) {
	biz := business.NewBusiness(logger, gormDb)

	return &Service{
		log:    logger,
		gormDb: gormDb,
		biz:    biz,
		client: client,
	}, nil
}
