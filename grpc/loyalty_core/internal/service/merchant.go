package service

import (
	"context"
	"loyalty_core/api"
)

func (s *Service) CreateMerchant(ctx context.Context, req *api.CreateMerchantRequest) (*api.CreateMerchantResponse, error) {
	err, id := s.biz.MerchantBusiness.ProcessCreateMerchant(ctx, req)
	if err != nil {
		return &api.CreateMerchantResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	account := &api.CreateAccountRequest{
		OwnerType: api.OwnerType_MERCHANT,
		OwnerId:   id,
		Point:     req.InitialPoints,
		Type:      api.AccountType_TOTAL,
	}

	_, err = s.client.CreateAccount(ctx, account)
	if err != nil {
		s.log.Fatal("Failed to create account")
		return &api.CreateMerchantResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	return &api.CreateMerchantResponse{
		Success: true,
		Message: "Merchant created successfully",
		Data: &api.CreateMerchantResponse_Data{
			MerchantId: id,
		},
	}, nil
}
