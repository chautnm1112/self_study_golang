package business

import (
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/repository"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/service/business/account"
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/service/business/transaction"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
