package service

import (
	"context"
	apiaccounting "github.com/chautnm1112/loyalty/loyalty_accounting/api"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
)

func (s *Service) CreateMerchant(ctx context.Context, req *api.CreateMerchantRequest) (*api.CreateMerchantResponse, error) {
	s.log.Named("SERVICE-MEMBER").Info("CreateMerchant start")

	err, id := s.biz.MerchantBusiness.ProcessCreateMerchant(ctx, req)
	if err != nil {
		return &api.CreateMerchantResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	account := &apiaccounting.CreateAccountRequest{
		OwnerType: apiaccounting.OwnerType_MERCHANT,
		OwnerId:   id,
		Point:     req.InitialPoints,
		Type:      apiaccounting.AccountType_TOTAL,
	}

	_, err = s.client.CreateAccount(ctx, account)
	if err != nil {
		return &api.CreateMerchantResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	s.log.Named("SERVICE-MEMBER").Info("CreateMerchant end")

	return &api.CreateMerchantResponse{
		Success: true,
		Message: "Merchant created successfully",
		Data: &api.CreateMerchantResponse_Data{
			MerchantId: id,
		},
	}, nil
}
