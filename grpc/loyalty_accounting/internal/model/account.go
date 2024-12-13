package model

import (
	"gorm.io/gorm"
	"loyalty_accounting/internal/enum"
	"time"
)

type Account struct {
	gorm.Model
	ID        uint32         `gorm:"primaryKey"`
	OwnerType enum.OwnerType `gorm:"not null;"`
	OwnerId   uint32
	Points    uint32
	Type      enum.AccountType `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Account) TableName() string {
	return "account"
}
