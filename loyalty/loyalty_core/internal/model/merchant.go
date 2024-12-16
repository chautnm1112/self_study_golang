package model

import (
	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model
	ID        uint32 `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null" json:"name"`
	NetworkId uint32 `json:"network_id"`
}

func (Merchant) TableName() string {
	return "merchants"
}
