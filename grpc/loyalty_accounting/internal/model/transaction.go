package model

import (
	"gorm.io/gorm"
	"time"
)
import "loyalty_accounting/internal/enum"

type Transaction struct {
	gorm.Model
	ID            uint32 `gorm:"primaryKey"`
	FromAccountId uint32
	ToAccountId   uint32
	Points        uint32
	Type          enum.TransactionType
	CreatedAt     time.Time
}

func (Transaction) TableName() string {
	return "transaction"
}
