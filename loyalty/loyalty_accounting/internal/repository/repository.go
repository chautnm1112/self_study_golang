package repository

import "gorm.io/gorm"

type Repository struct {
	AccountRepository
	TransactionRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		NewAccountRepository(db),
		NewTransactionRepository(db),
	}
}
