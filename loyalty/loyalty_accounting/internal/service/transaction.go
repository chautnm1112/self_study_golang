package service

import (
	"context"
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
)

func (s *Service) CreateTransaction(ctx context.Context, req *api.CreateTransactionRequest) (*api.CreateTransactionResponse, error) {
	err, id := s.biz.TransactionBusiness.ProcessCreateTransaction(ctx, req)
	if err != nil {
		return &api.CreateTransactionResponse{
			Success: false,
			Message: err.Error() + " happened",
			Data:    nil,
		}, nil
	}

	return &api.CreateTransactionResponse{
		Success: true,
		Message: "Transaction created successfully",
		Data: &api.CreateTransactionResponse_Data{
			TransactionId: id,
		},
	}, nil
}
