package model

import (
	"github.com/chautnm1112/loyalty/loyalty_accounting/internal/enum"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID        uint32         `gorm:"primaryKey"`
	OwnerType enum.OwnerType `gorm:"not null;"`
	OwnerId   uint32
	Points    uint32
	Type      enum.AccountType `gorm:"not null"`
}

func (Account) TableName() string {
	return "account"
}
