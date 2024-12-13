package model

import (
	"gorm.io/gorm"
	"time"
)

type Network struct {
	gorm.Model
	ID        uint32 `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null" json:"network_name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Network) TableName() string {
	return "networks"
}
