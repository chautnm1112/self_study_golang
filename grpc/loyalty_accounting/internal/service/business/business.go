package business

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"loyalty_accounting/internal/repository"
	"loyalty_accounting/internal/service/business/account"
	"loyalty_accounting/internal/service/business/transaction"
)

type Business struct {
	db                  *gorm.DB
	repository          *repository.Repository
	AccountBusiness     account.Business
	TransactionBusiness transaction.Business
}

func NewBusiness(
	log *zap.Logger,
	db *gorm.DB,
) *Business {
	repo := repository.NewRepository(db)
	return &Business{
		db:                  db,
		repository:          repo,
		AccountBusiness:     account.NewAccountBusiness(log, repo),
		TransactionBusiness: transaction.NewTransaction(log, repo),
	}
}
