package service

import (
	"context"
	"loyalty_accounting/api"
)

func (s *Service) CreateAccount(ctx context.Context, req *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {
	err, id := s.biz.AccountBusiness.ProcessCreateAccount(ctx, req)

	if err != nil {
		return &api.CreateAccountResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	return &api.CreateAccountResponse{
		Success: true,
		Message: "Points account created successfully",
		Data: &api.CreateAccountResponse_Data{
			AccountId: id,
		},
	}, nil
}
func (s *Service) EarnPoints(ctx context.Context, req *api.EarnPointsRequest) (*api.EarnPointsResponse, error) {
	err, transactionId := s.biz.AccountBusiness.ProcessEarnPoints(ctx, req)

	if err != nil {
		return &api.EarnPointsResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	return &api.EarnPointsResponse{
		Success: true,
		Message: "Earn points successfully",
		Data: &api.EarnPointsResponse_Data{
			TransactionId: transactionId,
		},
	}, nil
}
func (s *Service) RedeemPoints(ctx context.Context, req *api.RedeemPointsRequest) (*api.RedeemPointsResponse, error) {
	err, transactionId := s.biz.AccountBusiness.ProcessRedeemPoints(ctx, req)

	if err != nil {
		return &api.RedeemPointsResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	return &api.RedeemPointsResponse{
		Success: true,
		Message: "Redeem points successfully",
		Data: &api.RedeemPointsResponse_Data{
			TransactionId: transactionId,
		},
	}, nil
}
func (s *Service) RefundEarnedPoints(ctx context.Context, req *api.RefundEarnedPointsRequest) (*api.RefundEarnedPointsResponse, error) {
	err, transactionId := s.biz.AccountBusiness.ProcessRefundEarnedPoints(ctx, req)

	if err != nil {
		return &api.RefundEarnedPointsResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	return &api.RefundEarnedPointsResponse{
		Success: true,
		Message: "Refund earned points successfully",
		Data: &api.RefundEarnedPointsResponse_Data{
			TransactionId: transactionId,
		},
	}, nil
}
func (s *Service) RefundRedeemPoints(ctx context.Context, req *api.RefundRedeemPointsRequest) (*api.RefundRedeemPointsResponse, error) {
	err, transactionId := s.biz.AccountBusiness.ProcessRefundRedeemPoints(ctx, req)

	if err != nil {
		return &api.RefundRedeemPointsResponse{
			Success: false,
			Message: err.Error() + "happened",
			Data:    nil,
		}, nil
	}

	return &api.RefundRedeemPointsResponse{
		Success: true,
		Message: "Refund redeem points successfully",
		Data: &api.RefundRedeemPointsResponse_Data{
			TransactionId: transactionId,
		},
	}, nil
}
