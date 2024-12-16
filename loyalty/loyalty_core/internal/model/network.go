package model

import (
	"gorm.io/gorm"
)

type Network struct {
	gorm.Model
	ID   uint32 `gorm:"primaryKey"`
	Name string `gorm:"size:100;not null" json:"network_name"`
}

func (Network) TableName() string {
	return "networks"
}
