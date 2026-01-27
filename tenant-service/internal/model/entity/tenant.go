package entity

import (
	_ "encoding/json"
	"time"

	"gorm.io/gorm"
)

// Tenant 商户信息表
type Tenant struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"size:100;not null;index" json:"name"`
	Status    uint8          `gorm:"type:tinyint;default:1" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 软删除
}

func (Tenant) TableName() string {
	return "tenants"
}
