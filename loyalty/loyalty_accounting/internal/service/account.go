package service

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
)

func (s *Service) CreateAccount(ctx context.Context, req *api.CreateAccountRequest) (*api.CreateAccountResponse, error) {
	s.log.Named("SERVICE-ACCOUNT").Info("CreateAccount start")

	err, id := s.biz.AccountBusiness.ProcessCreateAccount(ctx, req)

	if err != nil {
		s.log.Named("SERVICE-ACCOUNT").Info("Create Account fail")

		return &api.CreateAccountResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	s.log.Named("SERVICE-ACCOUNT").Info("Create Account success")

	s.log.Named("SERVICE-ACCOUNT").Info("CreateAccount end")

	return &api.CreateAccountResponse{
		Success: true,
		Message: "Account created successfully",
		Data: &api.CreateAccountResponse_Data{
			AccountId: id,
		},
	}, nil
}

func (s *Service) UpdateAccount(ctx context.Context, req *api.UpdateAccountRequest) (*api.UpdateAccountResponse, error) {
	s.log.Named("SERVICE-ACCOUNT").Info("UpdateAccount start")

	err, id := s.biz.AccountBusiness.ProcessUpdateAccount(ctx, req)

	if err != nil {
		s.log.Named("SERVICE-ACCOUNT").Info("Update Account fail")

		return &api.UpdateAccountResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	s.log.Named("SERVICE-ACCOUNT").Info("Update Account success")

	s.log.Named("SERVICE-ACCOUNT").Info("UpdateAccount end")

	return &api.UpdateAccountResponse{
		Success: true,
		Message: "Account update successfully",
		Data: &api.UpdateAccountResponse_Data{
			AccountId: id,
		},
	}, nil
}
func (s *Service) EarnPoints(ctx context.Context, req *api.EarnPointsRequest) (*api.EarnPointsResponse, error) {
	s.log.Named("SERVICE-ACCOUNT").Info("EarnPoints start")

	err, transactionId := s.biz.AccountBusiness.ProcessEarnPoints(ctx, req)

	if err != nil {
		s.log.Named("SERVICE-ACCOUNT").Info("Earn Points fail")

		return &api.EarnPointsResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	s.log.Named("SERVICE-ACCOUNT").Info("Earn Points success")

	s.log.Named("SERVICE-ACCOUNT").Info("EarnPoints end")

	return &api.EarnPointsResponse{
		Success: true,
		Message: "Earn points successfully",
		Data: &api.EarnPointsResponse_Data{
			TransactionId: transactionId,
		},
	}, nil
}
func (s *Service) RedeemPoints(ctx context.Context, req *api.RedeemPointsRequest) (*api.RedeemPointsResponse, error) {
	s.log.Named("SERVICE-ACCOUNT").Info("RedeemPoints start")

	err, transactionId := s.biz.AccountBusiness.ProcessRedeemPoints(ctx, req)

	if err != nil {
		s.log.Named("SERVICE-ACCOUNT").Info("Redeem Points fail")

		return &api.RedeemPointsResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	s.log.Named("SERVICE-ACCOUNT").Info("Redeem Points success")

	s.log.Named("SERVICE-ACCOUNT").Info("RedeemPoints end")

	return &api.RedeemPointsResponse{
		Success: true,
		Message: "Redeem points successfully",
		Data: &api.RedeemPointsResponse_Data{
			TransactionId: transactionId,
		},
	}, nil
}
func (s *Service) RefundEarnedPoints(ctx context.Context, req *api.RefundEarnedPointsRequest) (*api.RefundEarnedPointsResponse, error) {
	s.log.Named("SERVICE-ACCOUNT").Info("RefundEarnedPoints start")

	err, transactionId := s.biz.AccountBusiness.ProcessRefundEarnedPoints(ctx, req)

	if err != nil {
		s.log.Named("SERVICE-ACCOUNT").Info("Refund Earned Points fail")

		return &api.RefundEarnedPointsResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	s.log.Named("SERVICE-ACCOUNT").Info("Refund Earned Points success")

	s.log.Named("SERVICE-ACCOUNT").Info("RefundEarnedPoints end")

	return &api.RefundEarnedPointsResponse{
		Success: true,
		Message: "Refund earned points successfully",
		Data: &api.RefundEarnedPointsResponse_Data{
			TransactionId: transactionId,
		},
	}, nil
}
func (s *Service) RefundRedeemPoints(ctx context.Context, req *api.RefundRedeemPointsRequest) (*api.RefundRedeemPointsResponse, error) {
	s.log.Named("SERVICE-ACCOUNT").Info("RefundRedeemPoints start")

	err, transactionId := s.biz.AccountBusiness.ProcessRefundRedeemPoints(ctx, req)

	if err != nil {
		s.log.Named("SERVICE-ACCOUNT").Info("Refund Redeem Points fail")

		return &api.RefundRedeemPointsResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	s.log.Named("SERVICE-ACCOUNT").Info("Refund Redeem Points success")

	s.log.Named("SERVICE-ACCOUNT").Info("RefundRedeemPoints end")

	return &api.RefundRedeemPointsResponse{
		Success: true,
		Message: "Refund redeem points successfully",
		Data: &api.RefundRedeemPointsResponse_Data{
			TransactionId: transactionId,
		},
	}, nil
}
