package model

import (
	"gorm.io/gorm"
	"time"
)

type Merchant struct {
	gorm.Model
	ID        uint32 `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null" json:"name"`
	NetworkId uint32 `json:"network_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Merchant) TableName() string {
	return "merchants"
}
