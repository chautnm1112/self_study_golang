package account

import (
	"context"
	"fmt"
	"github.com/chautnm1112/loyalty/loyalty_accounting/api"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/enum"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/model"
	"go.uber.org/zap"
)

func (b *business) ProcessCreateAccount(ctx context.Context, request *api.CreateAccountRequest) (error, uint32) {
	b.log.Info("ProcessCreateAccount start")

	account := b.mappingCreateAccount(request)

	err, accountId := b.repository.CreateAccount(ctx, nil, account)

	b.log.Info("CreateAccount done")

	if err != nil {
		b.log.Error("CreateAccount error", zap.Error(err))
		return err, 0
	}

	b.log.Info("ProcessCreateAccount end")

	return nil, accountId
}

func (b *business) ProcessUpdateAccount(ctx context.Context, request *api.UpdateAccountRequest) (error, uint32) {
	b.log.Info("ProcessUpdateAccount start")

	b.log.Info("api.UpdateAccountRequest", zap.Any("UpdateAccountRequest", request))

	account := b.mappingUpdateAccount(request)

	err, accountId := b.repository.UpdateAccount(ctx, nil, account)

	b.log.Info("UpdateAccount done")

	if err != nil {
		b.log.Error("UpdateAccount error", zap.Error(err))
		return err, 0
	}

	b.log.Info("Account updated", zap.Any("account", account))

	b.log.Info("ProcessUpdateAccount end")

	return err, accountId
}
func (b *business) ProcessEarnPoints(ctx context.Context, request *api.EarnPointsRequest) (error, uint32) {
	b.log.Info("ProcessEarnPoints start")

	b.log.Info("api.EarnPointsRequest", zap.Any("EarnPointsRequest", request))

	tx, err := b.repository.BeginTransaction(ctx)
	if err != nil {
		b.log.Error("BeginTransaction error", zap.Error(err))
		return err, 0
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	merchantAccount := &model.Account{}
	err = b.repository.GetAccountById(ctx, tx, request.GetMerchantAccountId(), merchantAccount)
	if err != nil {
		b.log.Error("GetAccountByOwnerIdAndOwnerType (Merchant) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("merchantAccount before", zap.Any("merchantAccount", merchantAccount))

	if merchantAccount.Points < request.GetPoints() {
		err := fmt.Errorf("insufficient points")
		b.log.Error("Merchant has insufficient points", zap.Error(err))
		return err, 0
	}

	merchantAccount.Points -= request.GetPoints()
	err, _ = b.repository.UpdateAccount(ctx, tx, merchantAccount)
	if err != nil {
		b.log.Error("UpdateAccount (Merchant) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("merchantAccount updated", zap.Any("merchantAccount", merchantAccount))

	transactionRequest := &model.Transaction{
		FromAccountId: request.MerchantAccountId,
		ToAccountId:   request.MemberAccountId,
		Points:        request.Points,
		Type:          enum.EARN_POINTS,
	}

	err, transactionId := b.repository.CreateTransaction(ctx, tx, transactionRequest)
	if err != nil {
		b.log.Error("CreateTransaction error", zap.Error(err))
		return err, 0
	}

	b.log.Info("transactionRequest", zap.Any("transactionRequest", transactionRequest))

	memberAccount := &model.Account{}
	err = b.repository.GetAccountById(ctx, tx, request.GetMemberAccountId(), memberAccount)
	if err != nil {
		b.log.Error("GetAccountByOwnerIdAndOwnerType (Member) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("memberAccount before", zap.Any("memberAccount", memberAccount))

	memberAccount.Points += request.GetPoints()
	err, _ = b.repository.UpdateAccount(ctx, tx, memberAccount)
	if err != nil {
		b.log.Error("UpdateAccount (Member) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("memberAccount updated", zap.Any("memberAccount", memberAccount))

	b.log.Info("ProcessEarnPoints end")

	return nil, transactionId
}

func (b *business) ProcessRedeemPoints(ctx context.Context, request *api.RedeemPointsRequest) (error, uint32) {
	b.log.Info("ProcessRedeemPoints start")

	b.log.Info("api.RedeemPointsRequest", zap.Any("RedeemPointsRequest", request))

	tx, err := b.repository.BeginTransaction(ctx)
	if err != nil {
		b.log.Error("BeginTransaction error", zap.Error(err))
		return err, 0
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	memberAccount := &model.Account{}
	err = b.repository.GetAccountById(ctx, tx, request.GetMemberAccountId(), memberAccount)
	if err != nil {
		b.log.Error("GetAccountByOwnerIdAndOwnerType (Member) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("memberAccount before", zap.Any("memberAccount", memberAccount))

	if memberAccount.Points < request.GetPoints() {
		err = fmt.Errorf("insufficient points")
		b.log.Error("Member has insufficient points", zap.Error(err))
		return err, 0
	}

	memberAccount.Points -= request.GetPoints()
	err, _ = b.repository.UpdateAccount(ctx, tx, memberAccount)
	if err != nil {
		b.log.Error("UpdateAccount (Member) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("memberAccount updated", zap.Any("memberAccount", memberAccount))

	networkAccount := &model.Account{}
	err = b.repository.GetAccountById(ctx, tx, request.GetNetworkAccountId(), networkAccount)
	if err != nil {
		b.log.Error("GetAccountByOwnerIdAndOwnerType (Network) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("networkAccount before", zap.Any("networkAccount", networkAccount))

	networkAccount.Points += request.GetPoints()
	err, _ = b.repository.UpdateAccount(ctx, tx, networkAccount)
	if err != nil {
		b.log.Error("UpdateAccount (Network) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("networkAccount updated", zap.Any("networkAccount", networkAccount))

	transactionRequest := &model.Transaction{
		FromAccountId: request.MemberAccountId,
		ToAccountId:   request.NetworkAccountId,
		Points:        request.Points,
		Type:          enum.REDEEM_POINTS,
	}

	err, transactionId := b.repository.CreateTransaction(ctx, tx, transactionRequest)
	if err != nil {
		b.log.Error("CreateTransaction error", zap.Error(err))
		return err, 0
	}

	b.log.Info("transactionRequest", zap.Any("transactionRequest", transactionRequest))

	b.log.Info("ProcessRedeemPoints end")

	return nil, transactionId
}

func (b *business) ProcessRefundEarnedPoints(ctx context.Context, request *api.RefundEarnedPointsRequest) (error, uint32) {
	b.log.Info("ProcessRefundEarnedPoints start")

	b.log.Info("api.RefundEarnedPointsRequest", zap.Any("RefundEarnedPointsRequest", request))

	tx, err := b.repository.BeginTransaction(ctx)
	if err != nil {
		b.log.Error("BeginTransaction error", zap.Error(err))
		return err, 0
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	merchantAccount := &model.Account{}
	err = b.repository.GetAccountById(ctx, tx, request.GetMerchantAccountId(), merchantAccount)
	if err != nil {
		b.log.Error("GetAccountByOwnerIdAndOwnerType (Merchant) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("merchantAccount before", zap.Any("merchantAccount", merchantAccount))

	merchantAccount.Points += request.GetPoints()
	err, _ = b.repository.UpdateAccount(ctx, tx, merchantAccount)
	if err != nil {
		b.log.Error("UpdateAccount (Merchant) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("merchantAccount updated", zap.Any("merchantAccount", merchantAccount))

	memberAccount := &model.Account{}
	err = b.repository.GetAccountById(ctx, tx, request.GetMemberAccountId(), memberAccount)
	if err != nil {
		b.log.Error("GetAccountByOwnerIdAndOwnerType (Member) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("memberAccount before", zap.Any("memberAccount", memberAccount))

	if memberAccount.Points < request.GetPoints() {
		err = fmt.Errorf("insufficient points")
		b.log.Error("Member has insufficient points", zap.Error(err))
		return err, 0
	}

	memberAccount.Points -= request.GetPoints()
	err, _ = b.repository.UpdateAccount(ctx, tx, memberAccount)
	if err != nil {
		b.log.Error("UpdateAccount (Member) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("memberAccount updated", zap.Any("memberAccount", memberAccount))

	transactionRequest := &model.Transaction{
		FromAccountId: request.GetMemberAccountId(),
		ToAccountId:   request.GetMerchantAccountId(),
		Points:        request.GetPoints(),
		Type:          enum.REFUND_EARNED_POINTS,
	}

	err, transactionId := b.repository.CreateTransaction(ctx, tx, transactionRequest)
	if err != nil {
		b.log.Error("CreateTransaction error", zap.Error(err))
		return err, 0
	}

	b.log.Info("transactionRequest created", zap.Any("transactionRequest", transactionRequest))

	b.log.Info("ProcessRefundEarnedPoints end")

	return nil, transactionId
}

func (b *business) ProcessRefundRedeemPoints(ctx context.Context, request *api.RefundRedeemPointsRequest) (error, uint32) {
	b.log.Info("ProcessRefundRedeemPoints start")

	b.log.Info("api.RefundRedeemPointsRequest", zap.Any("RefundRedeemPointsRequest", request))

	tx, err := b.repository.BeginTransaction(ctx)
	if err != nil {
		b.log.Error("BeginTransaction error", zap.Error(err))
		return err, 0
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	merchantAccount := &model.Account{}
	err = b.repository.GetAccountById(ctx, tx, request.GetMerchantAccountId(), merchantAccount)
	if err != nil {
		b.log.Error("GetAccountByOwnerIdAndOwnerType (Merchant) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("merchantAccount before", zap.Any("merchantAccount", merchantAccount))

	if merchantAccount.Points < request.GetPoints() {
		err = fmt.Errorf("insufficient points in MerchantAccount")
		b.log.Error("MerchantAccount has insufficient points", zap.Error(err))
		return err, 0
	}

	merchantAccount.Points -= request.GetPoints()
	err, _ = b.repository.UpdateAccount(ctx, tx, merchantAccount)
	if err != nil {
		b.log.Error("UpdateAccount (Merchant) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("merchantAccount updated", zap.Any("merchantAccount", merchantAccount))

	memberAccount := &model.Account{}
	err = b.repository.GetAccountById(ctx, tx, request.GetMemberAccountId(), memberAccount)
	if err != nil {
		b.log.Error("GetAccountByOwnerIdAndOwnerType (Member) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("memberAccount before", zap.Any("memberAccount", memberAccount))

	memberAccount.Points += request.GetPoints()
	err, _ = b.repository.UpdateAccount(ctx, tx, memberAccount)
	if err != nil {
		b.log.Error("UpdateAccount (Member) error", zap.Error(err))
		return err, 0
	}

	b.log.Info("memberAccount updated", zap.Any("memberAccount", memberAccount))

	transactionRequest := &model.Transaction{
		FromAccountId: request.GetMerchantAccountId(),
		ToAccountId:   request.GetMemberAccountId(),
		Points:        request.GetPoints(),
		Type:          enum.REFUND_REDEEMED_POINTS,
	}

	err, transactionId := b.repository.CreateTransaction(ctx, tx, transactionRequest)
	if err != nil {
		b.log.Error("CreateTransaction error", zap.Error(err))
		return err, 0
	}

	b.log.Info("transactionRequest created", zap.Any("transactionRequest", transactionRequest))

	b.log.Info("ProcessRefundRedeemPoints end")

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
	}
}

func (b *business) mappingUpdateAccount(request *api.UpdateAccountRequest) *model.Account {
	return &model.Account{
		OwnerType: mapOwnerType(request.OwnerType),
		OwnerId:   request.OwnerId,
		Points:    request.Point,
		Type:      mapAccountType(request.Type),
	}
}
