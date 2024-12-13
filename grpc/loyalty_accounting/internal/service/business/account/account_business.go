package account

import (
	"context"
	"loyalty_accounting/api"
	"loyalty_accounting/internal/enum"
	"loyalty_accounting/internal/model"
	"time"
)

func (b *business) ProcessCreateAccount(ctx context.Context, request *api.CreateAccountRequest) (error, uint32) {
	account := b.mappingCreateAccount(request)

	return b.repository.CreateAccount(ctx, account)
}

func (b *business) ProcessEarnPoints(ctx context.Context, request *api.EarnPointsRequest) (error, uint32) {
	transactionRequest := &model.Transaction{
		FromAccountId: request.MerchantAccountId,
		ToAccountId:   request.MemberAccountId,
		Points:        request.Points,
		Type:          enum.EARN_POINTS,
	}
	var err error
	var transactionId uint32

	err, transactionId = b.repository.CreateTransaction(ctx, transactionRequest)

	var merchantAccount *model.Account
	err = b.repository.GetAccountByOwnerIdAndOwnerType(ctx, request.GetMerchantAccountId(), enum.MERCHANT, merchantAccount)
	if err != nil {
		return err, 0
	}

	merchantAccount.Points -= request.GetPoints()

	err, _ = b.repository.UpdateAccount(ctx, merchantAccount)
	if err != nil {
		return err, 0
	}

	var memberAccount *model.Account
	err = b.repository.GetAccountByOwnerIdAndOwnerType(ctx, request.GetMemberAccountId(), enum.MEMBER, memberAccount)
	if err != nil {
		return err, 0
	}

	memberAccount.Points += request.GetPoints()

	err, _ = b.repository.UpdateAccount(ctx, memberAccount)
	if err != nil {
		return err, 0
	}

	return nil, transactionId
}

func (b *business) ProcessRedeemPoints(ctx context.Context, request *api.RedeemPointsRequest) (error, uint32) {
	transactionRequest := &model.Transaction{
		FromAccountId: request.MemberAccountId,
		ToAccountId:   request.NetworkAccountId,
		Points:        request.Points,
		Type:          enum.REDEEM_POINTS,
	}
	var err error
	var transactionId uint32

	err, transactionId = b.repository.CreateTransaction(ctx, transactionRequest)

	var networkAccount *model.Account
	err = b.repository.GetAccountByOwnerIdAndOwnerType(ctx, request.GetNetworkAccountId(), enum.NETWORK, networkAccount)
	if err != nil {
		return err, 0
	}

	networkAccount.Points += request.GetPoints()

	err, _ = b.repository.UpdateAccount(ctx, networkAccount)
	if err != nil {
		return err, 0
	}

	var memberAccount *model.Account
	err = b.repository.GetAccountByOwnerIdAndOwnerType(ctx, request.GetMemberAccountId(), enum.MEMBER, memberAccount)
	if err != nil {
		return err, 0
	}

	memberAccount.Points -= request.GetPoints()

	err, _ = b.repository.UpdateAccount(ctx, memberAccount)
	if err != nil {
		return err, 0
	}

	return nil, transactionId
}

func (b *business) ProcessRefundEarnedPoints(ctx context.Context, request *api.RefundEarnedPointsRequest) (error, uint32) {
	transactionRequest := &model.Transaction{
		FromAccountId: request.GetMemberAccountId(),
		ToAccountId:   request.GetMerchantAccountId(),
		Points:        request.Points,
		Type:          enum.REFUND_EARNED_POINTS,
	}
	var err error
	var transactionId uint32

	err, transactionId = b.repository.CreateTransaction(ctx, transactionRequest)
	if err != nil {
		return err, 0
	}

	var merchantAccount *model.Account
	err = b.repository.GetAccountByOwnerIdAndOwnerType(ctx, request.GetMerchantAccountId(), enum.MERCHANT, merchantAccount)
	if err != nil {
		return err, 0
	}

	merchantAccount.Points += request.GetPoints()

	err, _ = b.repository.UpdateAccount(ctx, merchantAccount)
	if err != nil {
		return err, 0
	}

	var memberAccount *model.Account
	err = b.repository.GetAccountByOwnerIdAndOwnerType(ctx, request.GetMemberAccountId(), enum.MEMBER, memberAccount)
	if err != nil {
		return err, 0
	}

	memberAccount.Points -= request.GetPoints()

	err, _ = b.repository.UpdateAccount(ctx, memberAccount)
	if err != nil {
		return err, 0
	}

	return nil, transactionId
}

func (b *business) ProcessRefundRedeemPoints(ctx context.Context, request *api.RefundRedeemPointsRequest) (error, uint32) {
	transactionRequest := &model.Transaction{
		FromAccountId: request.GetMerchantAccountId(),
		ToAccountId:   request.GetMemberAccountId(),
		Points:        request.Points,
		Type:          enum.REFUND_REDEEMED_POINTS,
	}
	var err error
	var transactionId uint32

	err, transactionId = b.repository.CreateTransaction(ctx, transactionRequest)
	if err != nil {
		return err, 0
	}

	var merchantAccount *model.Account
	err = b.repository.GetAccountByOwnerIdAndOwnerType(ctx, request.GetMerchantAccountId(), enum.MERCHANT, merchantAccount)
	if err != nil {
		return err, 0
	}

	merchantAccount.Points -= request.GetPoints()

	err, _ = b.repository.UpdateAccount(ctx, merchantAccount)
	if err != nil {
		return err, 0
	}

	var memberAccount *model.Account
	err = b.repository.GetAccountByOwnerIdAndOwnerType(ctx, request.GetMemberAccountId(), enum.MEMBER, memberAccount)
	if err != nil {
		return err, 0
	}

	memberAccount.Points += request.GetPoints()

	err, _ = b.repository.UpdateAccount(ctx, memberAccount)
	if err != nil {
		return err, 0
	}

	return nil, transactionId
}

func mapOwnerType(protoType api.OwnerType) enum.OwnerType {
	switch protoType {
	case api.OwnerType_NETWORK:
		return enum.NETWORK
	case api.OwnerType_MERCHANT:
		return enum.MERCHANT
	case api.OwnerType_MEMBER:
		return enum.MEMBER
	default:
		return enum.NETWORK
	}
}

func mapAccountType(protoType api.AccountType) enum.AccountType {
	switch protoType {
	case api.AccountType_TOTAL:
		return enum.TOTAL
	case api.AccountType_NETWORK_PROMOTION:
		return enum.NETWORK_PROMOTION
	case api.AccountType_NETWORK_REVOKE:
		return enum.NETWORK_REVOKE
	default:
		return ""
	}
}

func (b *business) mappingCreateAccount(request *api.CreateAccountRequest) *model.Account {
	return &model.Account{
		OwnerType: mapOwnerType(request.OwnerType),
		OwnerId:   request.OwnerId,
		Points:    request.Point,
		Type:      mapAccountType(request.Type),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (b *business) mappingUpdateAccount(request *api.UpdateAccountRequest) *model.Account {
	return &model.Account{
		OwnerType: enum.OwnerType(request.OwnerType),
		OwnerId:   request.OwnerId,
		Points:    request.Point,
		Type:      enum.AccountType(request.Type),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
