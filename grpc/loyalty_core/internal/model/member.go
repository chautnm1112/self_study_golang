package model

import (
	"gorm.io/gorm"
	"time"
)

type Member struct {
	gorm.Model
	ID        uint32 `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null" json:"name"`
	Email     string `gorm:"size:100" json:"email"`
	Phone     string `gorm:"size:100" json:"phone"`
	NetworkId uint32 `json:"network_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Member) TableName() string {
	return "members"
}
