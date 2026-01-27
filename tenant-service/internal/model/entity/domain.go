package entity

import (
	"time"

	"gorm.io/gorm"
)

// Domain 域名信息表
type Domain struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`        // 关联的商户ID
	Domain    string         `gorm:"size:255;not null;unique" json:"domain"` // 域名，与tenant组合唯一
	Status    uint8          `gorm:"type:tinyint;default:1" json:"status"`   // 域名状态
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 软删除
}

func (Domain) TableName() string {
	return "domains"
}
