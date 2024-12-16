package model

import (
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/enum"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID            uint32 `gorm:"primaryKey"`
	FromAccountId uint32
	ToAccountId   uint32
	Points        uint32
	Type          enum.TransactionType
}

func (Transaction) TableName() string {
	return "transaction"
}
