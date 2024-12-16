package service

import (
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
	apiLoyalty "github.com/chautnm1112/loyalty/loyalty_core/api"
	"github.com/chautnm1112/loyalty/loyalty_core/internal/service/business"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	log    *zap.Logger
	gormDb *gorm.DB
	biz    *business.Business
	apiLoyalty.UnimplementedLoyaltyCoreServiceServer
	client api.LoyaltyAccountingServiceClient
}

func (s *Service) mustEmbedUnimplementedLoyaltyCoreServiceServer() {
	//TODO implement me
	panic("implement me")
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
