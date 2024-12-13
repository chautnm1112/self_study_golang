package service

import (
	"context"
	"loyalty_core/api"
)

func (s *Service) CreateNetwork(ctx context.Context, req *api.CreateNetworkRequest) (*api.CreateNetworkResponse, error) {
	err, id := s.biz.NetworkBusiness.ProcessCreateNetwork(ctx, req)
	if err != nil {
		return &api.CreateNetworkResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	account := &api.CreateAccountRequest{
		OwnerType: api.OwnerType_NETWORK,
		OwnerId:   id,
		Point:     req.InitialPoints,
		Type:      api.AccountType_TOTAL,
	}

	_, err = s.client.CreateAccount(ctx, account)
	if err != nil {
		s.log.Fatal("Failed to creat account")
		return &api.CreateNetworkResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	account = &api.CreateAccountRequest{
		OwnerType: api.OwnerType_NETWORK,
		OwnerId:   id,
		Point:     req.InitialPoints,
		Type:      api.AccountType_NETWORK_PROMOTION,
	}

	_, err = s.client.CreateAccount(ctx, account)

	return &api.CreateNetworkResponse{
		Success: true,
		Message: "Network created successfully",
		Data: &api.CreateNetworkResponse_Data{
			NetworkId: id,
		},
	}, nil
}
