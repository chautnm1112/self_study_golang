package service

import (
	"context"
	apiaccounting "github.com/chautnm1112/loyalty/loyalty_accounting/api"
	"github.com/chautnm1112/loyalty/loyalty_core/api"
)

func (s *Service) CreateNetwork(ctx context.Context, req *api.CreateNetworkRequest) (*api.CreateNetworkResponse, error) {
	s.log.Named("SERVICE-MEMBER").Info("CreateNetwork start")

	err, id := s.biz.NetworkBusiness.ProcessCreateNetwork(ctx, req)
	if err != nil {
		return &api.CreateNetworkResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	account := &apiaccounting.CreateAccountRequest{
		OwnerType: apiaccounting.OwnerType_NETWORK,
		OwnerId:   id,
		Point:     req.InitialPoints,
		Type:      apiaccounting.AccountType_TOTAL,
	}

	_, err = s.client.CreateAccount(ctx, account)
	if err != nil {
		return &api.CreateNetworkResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	account = &apiaccounting.CreateAccountRequest{
		OwnerType: apiaccounting.OwnerType_NETWORK,
		OwnerId:   id,
		Point:     req.InitialPoints,
		Type:      apiaccounting.AccountType_NETWORK_PROMOTION,
	}

	_, err = s.client.CreateAccount(ctx, account)

	s.log.Named("SERVICE-MEMBER").Info("CreateNetwork end")

	return &api.CreateNetworkResponse{
		Success: true,
		Message: "Network created successfully",
		Data: &api.CreateNetworkResponse_Data{
			NetworkId: id,
		},
	}, nil
}
